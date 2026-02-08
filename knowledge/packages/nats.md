# NATS Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a NATS client package. It uses the `nats-io/nats.go` library and follows the standard module principles in `guide/module.md`.

### Architectural Role

*   **Location**: `pkg/adapter/nats`
*   **Purpose**: To provide a concrete implementation for publishing, subscribing, and requesting on the NATS messaging system.

### Generation Prompt for LLM

When asked to create a NATS client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/adapter/nats/fx.go`**

```go
package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.uber.org/fx"

	// e.g., "your/project/lib/messaging"
)

var Module = fx.Module("nats-adapter",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// fx.As(new(messaging.Broker)),
		),
		ConfigRegister,
	),
)

type Config struct {
	URL string `yaml:"url"`
}

func ConfigRegister() *Config {
	return &Config{
		URL: nats.DefaultURL,
	}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

type Client struct {
	*nats.Conn
}

func NewClient(p Param) (*Client, error) {
	nc, err := nats.Connect(p.Config.URL)
	if err != nil {
		return nil, err
	}

	client := &Client{Conn: nc}

	p.Lifecycle.Append(fx.Hook{
		OnStart: nil, // Connection is already established.
		OnStop: func(ctx context.Context) error {
			client.Conn.Drain()
			return nil
		},
	})

	return client, nil
}
```
