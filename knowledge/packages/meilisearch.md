# Meilisearch Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a Meilisearch client package. It uses the `meilisearch/meilisearch-go` library and follows the standard module principles in `guide/module.md`.

### Architectural Role

*   **Location**: `pkg/client/meilisearch`
*   **Purpose**: To provide a concrete implementation of a `lib` interface for fast, relevant search, encapsulating all interactions with a Meilisearch instance.

### Generation Prompt for LLM

When asked to create a Meilisearch client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/client/meilisearch/fx.go`**

```go
package meilisearch

import (
	"context"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/fx"

	// e.g., "your/project/lib/search"
)

var Module = fx.Module("meilisearch-client",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// fx.As(new(search.Engine)),
		),
		ConfigRegister,
	),
)

type Config struct {
	Host   string `yaml:"host"`
	APIKey string `yaml:"apiKey"`
}

func ConfigRegister() *Config {
	return &Config{
		Host: "http://localhost:7700",
	}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

type Client struct {
	*meilisearch.Client
}

func NewClient(p Param) *Client {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   p.Config.Host,
		APIKey: p.Config.APIKey,
	})

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := client.Health()
			return err
		},
		OnStop: nil,
	})

	return &Client{Client: client}
}
```
