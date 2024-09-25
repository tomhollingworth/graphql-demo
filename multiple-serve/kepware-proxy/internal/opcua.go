package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/thollingworth/graphql-demo/multiple-serve/kepware-proxy/domain"
)

type OpcuaClient struct {
	client *opcua.Client
}

func NewOpcuaClient(ctx context.Context, endpoint string) (*OpcuaClient, error) {

	c, err := opcua.NewClient(endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err != nil {
		return nil, err
	}
	if err := c.Connect(ctx); err != nil {
		return nil, err
	}

	// Close the connection when the context is done
	go func(ctx context.Context, c *opcua.Client) {
		<-ctx.Done()
		slog.Info("closing opcua client")
		c.Close(ctx)
	}(ctx, c)

	return &OpcuaClient{client: c}, nil
}

func (o *OpcuaClient) Read(ctx context.Context, address string) (*domain.Tag, error) {
	if o.client == nil {
		return nil, fmt.Errorf("opcua client not connected")
	}

	id, err := ua.ParseNodeID(address)
	if err != nil {
		log.Fatalf("invalid node id: %v", err)
		return nil, err
	}

	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead: []*ua.ReadValueID{
			{NodeID: id},
		},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	var resp *ua.ReadResponse
	for {
		resp, err = o.client.Read(ctx, req)
		if err == nil {
			break
		}

		// Following switch contains known errors that can be retried by the user.
		// Best practice is to do it on read operations.
		switch {
		case err == io.EOF && o.client.State() != opcua.Closed:
			// has to be retried unless user closed the connection
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionIDInvalid):
			// Session is not activated has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionNotActivated):
			// Session is invalid has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSecureChannelIDInvalid):
			// secure channel will be recreated internally.
			time.After(1 * time.Second)
			continue

		default:
			log.Fatalf("Read failed: %s", err)
		}
	}

	if resp != nil && resp.Results[0].Status != ua.StatusOK {
		err := fmt.Errorf("Status not OK: %v", resp.Results[0].Status)
		slog.Error(err.Error())
		return nil, err
	}

	result := resp.Results[0]

	dt := domain.UATypetoDataType(result.Value.Type())

	return &domain.Tag{
		Address:   address,
		Value:     result.Value.Value(),
		Datatype:  dt,
		Timestamp: result.ServerTimestamp,
	}, nil
}
