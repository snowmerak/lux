# ClickHouse Client Package (`pkg`) Generation Guide

This guide provides a detailed blueprint for creating a ClickHouse client package within the `pkg/client/` directory. It follows the standard module principles in `guide/module.md` and uses the `ClickHouse/clickhouse-go/v2` library.

### Architectural Role

*   **Location**: `pkg/client/clickhouse`
*   **Purpose**: To provide a concrete implementation of one or more `lib` interfaces that require a fast, column-oriented analytics database. This module encapsulates all interactions with the ClickHouse server.

### Generation Prompt for LLM

When asked to create a ClickHouse client package, follow these steps to generate the required files.

#### 1. Create the Directory Structure

Create the directory `pkg/client/clickhouse/`.

#### 2. Create the `fx.go` Module Definition

This file wires the ClickHouse client into the application's dependency injection system.

**File: `pkg/client/clickhouse/fx.go`**

```go
package clickhouse

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/fx"

	// Import the specific lib interface(s) this client will implement.
	// e.g., "your/project/lib/analytics"
)

// Module exports the ClickHouse client functionality to the Fx application.
var Module = fx.Module("clickhouse-client",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// fx.As(new(analytics.Engine)), // Example: cast to an interface
		),
		ConfigRegister,
	),
)

// Config holds the configuration for the ClickHouse client.
type Config struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// ConfigRegister loads the ClickHouse configuration.
func ConfigRegister() *Config {
	return &Config{
		Addr:     "localhost:9000",
		User:     "default",
		Password: "",
		Database: "default",
	}
}

// Param groups the dependencies for the ClickHouse client constructor.
type Param struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
}

// Client is the concrete implementation for ClickHouse.
type Client struct {
	clickhouse.Conn
}

// NewClient is the constructor for the ClickHouse client.
func NewClient(p Param) (*Client, error) {
	// The compile-time interface check is mandatory.
	// var _ analytics.Engine = (*Client)(nil)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{p.Config.Addr},
		Auth: clickhouse.Auth{
			Database: p.Config.Database,
			Username: p.Config.User,
			Password: p.Config.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	client := &Client{Conn: conn}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Conn.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return client.Conn.Close()
		},
	})

	return client, nil
}
```

#### 3. Implement the Interface Methods

Implement the `lib` interface methods in a separate file.

**File: `pkg/client/clickhouse/analytics.go`**

```go
package clickhouse

import "context"

// Example implementation of a method.
/*
func (c *Client) LogEvent(ctx context.Context, event string) error {
    return c.Conn.Exec(ctx, "INSERT INTO events (name) VALUES (?)", event)
}
*/
```
