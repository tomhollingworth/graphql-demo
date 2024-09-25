package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/thollingworth/graphql-demo/multiple-serve/kepware-proxy/graph"
	"github.com/thollingworth/graphql-demo/multiple-serve/kepware-proxy/internal"
)

const defaultPort = "8081"

var (
	connectionString = "opc.tcp://127.0.0.1:49320"
)

func main() {
	ctx := context.Background()
	defer ctx.Done()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// check if the endpoint is set in the environment
	endpoint := os.Getenv("OPCUA_ENDPOINT")
	if endpoint != "" {
		connectionString = endpoint
	}

	client, err := internal.NewOpcuaClient(ctx, connectionString)
	if err != nil {
		slog.Error("failed to create opc-ua client: ", err)
		return
	}

	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					Client: client,
				},
			},
		),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
