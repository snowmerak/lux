---
description: Guide for generating an Elasticsearch client Fx module using go-elasticsearch.
tags: [go, elasticsearch, search, analytics, client, fx, module]
---

# Elasticsearch Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating an Elasticsearch client package. It uses the `github.com/elastic/go-elasticsearch/v8` library and follows the standard module principles in `guide/module.md`.

### Architectural Role

*   **Location**: `pkg/client/elasticsearch`
*   **Purpose**: To provide a concrete implementation of a `lib` interface for search and analytics, encapsulating all interactions with an Elasticsearch cluster.

### Generation Prompt for LLM

When asked to create an Elasticsearch client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/client/elasticsearch/fx.go`**

```go
package elasticsearch

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/fx"

	// e.g., "your/project/lib/search"
)

var Module = fx.Module("elasticsearch-client",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// fx.As(new(search.Engine)),
		),
		ConfigRegister,
	),
)

type Config struct {
	Addresses []string `yaml:"addresses"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
}

func ConfigRegister() *Config {
	return &Config{
		Addresses: []string{"http://localhost:9200"},
	}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

type Client struct {
	*elasticsearch.Client
}

func NewClient(p Param) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: p.Config.Addresses,
		Username:  p.Config.Username,
		Password:  p.Config.Password,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	client := &Client{Client: es}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := client.Info(client.Info.WithContext(ctx))
			return err
		},
		OnStop: nil, // The client does not have an explicit close method.
	})

	return client, nil
}
```
