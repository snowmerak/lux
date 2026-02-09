---
description: Guide for generating an OpenSearch client Fx module for search and analytics.
tags: [go, opensearch, search, fx, module]
---

# OpenSearch Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating an OpenSearch client package. It uses the `opensearch-project/opensearch-go/v2` library and follows the standard module principles in `guide/module.md`.

### Architectural Role

*   **Location**: `pkg/client/opensearch`
*   **Purpose**: To provide a concrete implementation of a `lib` interface for search and analytics, encapsulating all interactions with an OpenSearch cluster.

### Generation Prompt for LLM

When asked to create an OpenSearch client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/client/opensearch/fx.go`**

```go
package opensearch

import (
	"context"

	"github.com/opensearch-project/opensearch-go/v2"
	"go.uber.org/fx"

	// e.g., "your/project/lib/search"
)

var Module = fx.Module("opensearch-client",
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
		Username:  "admin",
		Password:  "admin",
	}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

type Client struct {
	*opensearch.Client
}

func NewClient(p Param) (*Client, error) {
	cfg := opensearch.Config{
		Addresses: p.Config.Addresses,
		Username:  p.Config.Username,
		Password:  p.Config.Password,
	}
	os, err := opensearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	client := &Client{Client: os}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := client.Info()
			return err
		},
		OnStop: nil,
	})

	return client, nil
}
```
