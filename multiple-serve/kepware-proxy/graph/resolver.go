package graph

import "github.com/thollingworth/graphql-demo/multiple-serve/kepware-proxy/internal"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Client *internal.OpcuaClient
}
