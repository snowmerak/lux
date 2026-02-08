# Cassandra Client Package (`pkg`) Generation Guide

This guide provides a detailed blueprint for creating a Cassandra client package within the `pkg/client/` directory. It follows the standard module principles defined in `guide/module.md` and uses the `gocql/gocql` library.

### Architectural Role

*   **Location**: `pkg/client/cassandra`
*   **Purpose**: To provide a concrete implementation of one or more `lib` interfaces that require a distributed NoSQL database. This module encapsulates all interactions with the Cassandra cluster.

### Generation Prompt for LLM

When asked to create a Cassandra client package, follow these steps to generate the required files.

#### 1. Create the Directory Structure

Create the directory `pkg/client/cassandra/`.

#### 2. Create the `fx.go` Module Definition

This file wires the Cassandra client into the application's dependency injection system.

**File: `pkg/client/cassandra/fx.go`**

```go
package cassandra

import (
	"context"

	"github.com/gocql/gocql"
	"go.uber.org/fx"

	// Import the specific lib interface(s) this client will implement.
	// e.g., "your/project/lib/datastore"
)

// Module exports the Cassandra client functionality to the Fx application.
var Module = fx.Module("cassandra-client",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// fx.As(new(datastore.Store)), // Example: cast to an interface
		),
		ConfigRegister,
	),
)

// Config holds the configuration for the Cassandra client.
type Config struct {
	Hosts    []string `yaml:"hosts"`
	Keyspace string   `yaml:"keyspace"`
}

// ConfigRegister loads the Cassandra configuration.
func ConfigRegister() *Config {
	return &Config{
		Hosts:    []string{"127.0.0.1:9042"},
		Keyspace: "system",
	}
}

// Param groups the dependencies for the Cassandra client constructor.
type Param struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
}

// Client is the concrete implementation for Cassandra.
type Client struct {
	*gocql.Session
}

// NewClient is the constructor for the Cassandra client.
func NewClient(p Param) (*Client, error) {
	// The compile-time interface check is mandatory.
	// var _ datastore.Store = (*Client)(nil)

	cluster := gocql.NewCluster(p.Config.Hosts...)
	cluster.Keyspace = p.Config.Keyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	client := &Client{Session: session}

	p.Lifecycle.Append(fx.Hook{
		OnStart: nil, // The session is already created and ready.
		OnStop: func(ctx context.Context) error {
			// Gracefully close the Cassandra session.
			client.Session.Close()
			return nil
		},
	})

	return client, nil
}
```

#### 3. Implement the Interface Methods

Implement the `lib` interface methods in a separate file.

**File: `pkg/client/cassandra/store.go`**

```go
package cassandra

import "context"

// Example implementation of a method.
/*
func (c *Client) GetItem(ctx context.Context, id string) (string, error) {
    var result string
    if err := c.Session.Query(`SELECT data FROM items WHERE id = ?`, id).WithContext(ctx).Consistency(gocql.One).Scan(&result); err != nil {
        return "", err
    }
    return result, nil
}
*/
```
