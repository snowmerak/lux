---
description: Guide for generating a DuckDB client Fx module using go-duckdb.
tags: [go, duckdb, analytics, olap, database, client, fx, module]
---

# DuckDB Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a DuckDB client package within the `pkg/client/` directory. It uses the `github.com/marcboeker/go-duckdb` driver and follows the standard module principles in `guide/module.md`.

### Architectural Role

*   **Location**: `pkg/client/duckdb`
*   **Purpose**: To provide a concrete implementation of a `lib` interface that requires an in-process analytical database. This module encapsulates all interactions with the DuckDB database file.

### Generation Prompt for LLM

When asked to create a DuckDB client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/client/duckdb/fx.go`**

```go
package duckdb

import (
	"context"
	"database/sql"

	_ "github.com/marcboeker/go-duckdb"
	"go.uber.org/fx"

	// e.g., "your/project/lib/analytics"
)

var Module = fx.Module("duckdb-client",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// fx.As(new(analytics.Engine)),
		),
		ConfigRegister,
	),
)

type Config struct {
	DSN string `yaml:"dsn"`
}

func ConfigRegister() *Config {
	return &Config{
		DSN: "./local.db",
	}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

type Client struct {
	*sql.DB
}

func NewClient(p Param) (*Client, error) {
	db, err := sql.Open("duckdb", p.Config.DSN)
	if err != nil {
		return nil, err
	}

	client := &Client{DB: db}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.DB.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return client.DB.Close()
		},
	})

	return client, nil
}
```
