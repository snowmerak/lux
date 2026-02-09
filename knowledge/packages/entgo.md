---
description: Guide for generating an Entgo ORM Fx module.
tags: [go, entgo, orm, database, sql, ent, fx, module]
---

# Entgo ORM Package (`pkg`) Generation Guide

This guide provides a blueprint for creating an `entgo` ORM package. It follows the standard module principles and depends on a database client module (e.g., `pkg/client/postgres`).

### Architectural Role

*   **Location**: `pkg/orm/ent`
*   **Purpose**: To provide a type-safe data access layer using the `entgo` ORM. This module defines the database schema in Go code, runs migrations, and provides a client for all data operations.

### Generation Prompt for LLM

When asked to create an `entgo` package, follow these steps.

#### 1. Define the Schema

Create schema files in `pkg/orm/ent/schema/`. For example, `user.go`:

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}
```

#### 2. Generate the Ent Code

Run `go generate ./...` in the `pkg/orm/ent` directory. This will generate the ORM client and entity code.

#### 3. Create the `fx.go` Module Definition

**File: `pkg/orm/ent/fx.go`**

```go
package ent

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"go.uber.org/fx"

	// Import a database driver, e.g., from a postgres client module
	_ "github.com/jackc/pgx/v5/stdlib"
)

var Module = fx.Module("ent-orm",
	fx.Provide(
		NewClient,
		ConfigRegister,
	),
)

type Config struct {
	Migrate bool `yaml:"migrate"`
}

func ConfigRegister() *Config {
	return &Config{Migrate: true}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
	DB        *sql.DB // Depends on a *sql.DB from another module
}

func NewClient(p Param) (*Client, error) {
	drv := entsql.OpenDB(dialect.Postgres, p.DB)
	client := NewClient(Driver(drv))

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if p.Config.Migrate {
				return client.Schema.Create(ctx)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
```
