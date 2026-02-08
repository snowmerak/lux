# Qdrant Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a Qdrant client package. It uses the `qdrant/go-client` library and follows the standard module principles in `guide/module.md`.

### Architectural Role

*   **Location**: `pkg/client/qdrant`
*   **Purpose**: To provide a concrete implementation of a `lib` interface for vector search, encapsulating all interactions with a Qdrant instance.

### Generation Prompt for LLM

When asked to create a Qdrant client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/client/qdrant/fx.go`**

```go
package qdrant

import (
	"context"

	qdrant "github.com/qdrant/go-client/qdrant"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	// e.g., "your/project/lib/vectorsearch"
)

var Module = fx.Module("qdrant-client",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// fx.As(new(vectorsearch.Engine)),
		),
		ConfigRegister,
	),
)

type Config struct {
	Host string `yaml:"host"`
	Port uint `yaml:"port"`
}

func ConfigRegister() *Config {
	return &Config{
		Host: "localhost",
		Port: 6333,
	}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

// Client wraps the Qdrant points client.
type Client struct {
	qdrant.PointsClient
	conn *grpc.ClientConn
}

func NewClient(p Param) (*Client, error) {
	conn, err := grpc.Dial(p.Config.Host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := &Client{
		PointsClient: qdrant.NewPointsClient(conn),
		conn:         conn,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			return client.conn.Close()
		},
	})

	return client, nil
}
```
