# Valkey Client Package (`pkg`) Generation Guide

This guide provides a detailed blueprint for creating a Valkey client package within the `pkg/client/` directory. It follows the standard module principles defined in `guide/module.md` and utilizes the official `valkey-io/valkey-go` library.

### Architectural Role

*   **Location**: `pkg/client/valkey`
*   **Purpose**: To provide a concrete implementation of one or more `lib` interfaces that require a key-value store for caching, session management, etc. This module encapsulates all interactions with the Valkey server.

### Generation Prompt for LLM

When asked to create a Valkey client package, follow these steps to generate the required files.

#### 1. Create the Directory Structure

Create the directory `pkg/client/valkey/`.

#### 2. Create the `fx.go` Module Definition

This file wires the Valkey client into the application's dependency injection system. It provides the client constructor and its configuration.

**File: `pkg/client/valkey/fx.go`**

```go
package valkey

import (
	"context"

	"github.com/valkey-io/valkey-go"
	"go.uber.org/fx"

	// Import the specific lib interface(s) this client will implement.
	// For example:
	// "your/project/lib/cache"
)

// Module exports the Valkey client functionality to the Fx application.
var Module = fx.Module("valkey-client",
	fx.Provide(
		// Provide the constructor for the Valkey client.
		// It is annotated to be cast to the specific interface it implements.
		fx.Annotate(
			NewClient,
			// fx.As(new(cache.Cache)), // Example: cast to a cache interface
		),
		// Provide the configuration constructor.
		ConfigRegister,
	),
)

// Config holds the configuration for the Valkey client.
type Config struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// ConfigRegister loads the Valkey configuration.
func ConfigRegister() *Config {
	// In a real application, load from a file or environment variables.
	return &Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

// Param groups the dependencies for the Valkey client constructor.
type Param struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
}

// Client is the concrete implementation for Valkey.
type Client struct {
	valkey.Client
}

// NewClient is the constructor for the Valkey client.
func NewClient(p Param) (Client, error) {
	// The compile-time interface check is mandatory.
	// var _ cache.Cache = (*Client)(nil)

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{p.Config.Addr},
		Password:    p.Config.Password,
		SelectDB:    p.Config.DB,
	})
	if err != nil {
		return nil, err
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// The valkey-go client automatically handles connections, so no explicit
			// action is required on start. The connection is established on the first command.
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Gracefully close the Valkey connection and release resources.
			client.Close()
			return nil
		},
	})

	return client, nil
}
```

#### 3. Implement the Interface Methods

If the `Client` struct implements a `lib` interface, the methods for that interface should be implemented in a separate file (e.g., `cache.go`).

**File: `pkg/client/valkey/cache.go`**

```go
package valkey

import "context"

// Example implementation of a Set method for a cache.Cache interface.
/*
func (c *Client) Set(ctx context.Context, key string, value string) error {
    return c.Do(ctx, c.B().Set().Key(key).Value(value).Build()).Error()
}
*/
```
