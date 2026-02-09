---
description: Guide for generating a Microsoft SQL Server client Fx module.
tags: [go, mssql, sql, database, fx, module]
---

# MSSQL Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a Microsoft SQL Server client package. It uses the `database/sql` package with the `microsoft/go-mssqldb` driver and follows the standard module principles.

### Architectural Role

*   **Location**: `pkg/client/mssql`
*   **Purpose**: To provide a concrete implementation of a `lib` interface that requires a connection to a SQL Server database.

### Generation Prompt for LLM

When asked to create an MSSQL client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/client/mssql/fx.go`**

```go
package mssql

import (
	"context"
	"database/sql"

	_ "github.com/microsoft/go-mssqldb"
	"go.uber.org/fx"

	// e.g., "your/project/lib/repository"
)

var Module = fx.Module("mssql-client",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// This typically provides a *sql.DB to be used by repository implementations
		),
		ConfigRegister,
	),
)

type Config struct {
	DSN string `yaml:"dsn"`
}

func ConfigRegister() *Config {
	return &Config{
		DSN: "sqlserver://username:password@localhost:1433?database=master",
	}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

// NewClient provides a *sql.DB connection pool.
func NewClient(p Param) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", p.Config.DSN)
	if err != nil {
		return nil, err
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return db.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}
```
