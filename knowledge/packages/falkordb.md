---
description: Guide for generating a FalkorDB client Fx module using falkordb-go.
tags: [go, falkordb, graph, database, client, fx, module]
---

# FalkorDB Client Package (`pkg`) Generation Guide

This guide provides a detailed blueprint for creating a FalkorDB client package within the `pkg/client/` directory. It follows the standard module principles and uses the `snowmerak/falkordb-go` library.

### Architectural Role

*   **Location**: `pkg/client/falkordb`
*   **Purpose**: To provide a concrete implementation of one or more `lib` interfaces (e.g., repository or graph database access) using FalkorDB.

### Generation Prompt for LLM

When asked to create a FalkorDB client package, follow these steps to generate the required files.

#### 1. Create the Directory Structure

Create the directory `pkg/client/falkordb/`.

#### 2. Create the `fx.go` Module Definition

This file wires the FalkorDB client into the application's dependency injection system.

**File: `pkg/client/falkordb/fx.go`**

```go
package falkordb

import (
	"github.com/snowmerak/falkordb-go"
	"go.uber.org/fx"
)

// Module exports the FalkorDB client functionality to the Fx application.
var Module = fx.Module("falkordb-client",
	fx.Provide(
		NewClient,
		ConfigRegister,
	),
)
```

#### 3. Create the `config.go` Configuration

**File: `pkg/client/falkordb/config.go`**

```go
package falkordb

import (
	"fmt"
	"os"
)

type Config struct {
	URL string // e.g., falkor://0.0.0.0:6379
}

func ConfigRegister() (*Config, error) {
	url := os.Getenv("FALKORDB_URL")
	if url == "" {
		return nil, fmt.Errorf("FALKORDB_URL environment variable is required")
	}
	return &Config{URL: url}, nil
}
```

#### 4. Create the `client.go` Implementation

This file manages the lifecycle of the FalkorDB client.

**File: `pkg/client/falkordb/client.go`**

```go
package falkordb

import (
	"context"

	"github.com/snowmerak/falkordb-go"
	"go.uber.org/fx"
)

type Client struct {
	*falkordb.FalkorDB
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

func NewClient(p Param) (*Client, error) {
	db, err := falkordb.FromURL(p.Config.URL)
	if err != nil {
		return nil, err
	}

	p.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return &Client{FalkorDB: db}, nil
}
```

#### 5. Implement a specific business interface (Example)

If implementing a specific interface from `lib/repository/social.go`:

**File: `pkg/client/falkordb/social_repository.go`**

```go
package falkordb

import (
	"context"
	"your/project/lib/repository"
	"github.com/snowmerak/falkordb-go/graph"
)

var _ repository.SocialRepository = (*socialRepository)(nil)

type socialRepository struct {
	client *Client
}

func NewSocialRepository(client *Client) repository.SocialRepository {
	return &socialRepository{client: client}
}

func (r *socialRepository) CreatePerson(ctx context.Context, name string, age int) error {
	g := r.client.SelectGraph("social")
	_, err := g.Query("CREATE (:Person {name:$name, age:$age})", map[string]interface{}{
		"name": name,
		"age":  age,
	}, nil)
	return err
}
```
