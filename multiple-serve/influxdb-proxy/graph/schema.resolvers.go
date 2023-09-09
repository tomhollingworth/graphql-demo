package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.37

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/tomhollingworth/graphql-demo/multiple-serve/influxdb-proxy/domain"
)

// CreateHistory is the resolver for the createHistory field.
func (r *mutationResolver) CreateHistory(ctx context.Context, input domain.NewHistory) (*domain.History, error) {
	history := domain.History{
		Timestamp: input.Timestamp,
		Value:     input.Value,
		EquipmentProperty: &domain.EquipmentProperty{
			ID: input.PropertyID,
		},
		Datatype: input.Datatype,
	}

	writeAPI := client.WriteAPIBlocking(org, bucket)
	// Create point using full params constructor
	p := influxdb2.NewPoint("values",
		map[string]string{"property": history.EquipmentProperty.ID},
		map[string]interface{}{"value": history.Value, "datatype": history.Datatype},
		history.Timestamp)
	// Write point immediately
	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		return nil, err
	}
	return &history, nil
}

// History is the resolver for the history field.
func (r *queryResolver) History(ctx context.Context, filter domain.FilterHistory) ([]*domain.History, error) {
	queryAPI := client.QueryAPI(org)
	// Get QueryTableResult

	q := `from(bucket:"` + bucket + `")
`
	if filter.Timestamp == nil {
		q += ` |> range(start: -1h)
`
	} else {
		q += ` |> range(start: ` + filter.Timestamp.Min.Format(time.RFC3339) + `, stop: ` + filter.Timestamp.Max.Format(time.RFC3339) + `)
`
	}
	q += ` |> filter(fn: (r) => r._measurement == "values")
`
	if filter.PropertyID != nil {
		q += ` |> filter(fn: (r) => r.property == "` + *filter.PropertyID + `")
`
	}
	q += `|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
`

	fmt.Println(q)
	result, err := queryAPI.Query(context.Background(), q)
	history := make([]*domain.History, 0)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			dt, ok := result.Record().ValueByKey("datatype").(string)
			if !ok {
				fmt.Printf("failed to parse datatype: %v\n", result.Record().ValueByKey("datatype"))
			}
			v := result.Record().ValueByKey("value")
			h := domain.History{
				EquipmentProperty: &domain.EquipmentProperty{
					ID: result.Record().ValueByKey("property").(string),
				},
				Datatype:  domain.DataType(dt),
				Value:     fmt.Sprintf("%v", v),
				Timestamp: result.Record().Time(),
			}
			history = append(history, &h)
		}
		// Check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		return nil, err
	}
	// Ensures background processes finishes
	client.Close()

	return history, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
