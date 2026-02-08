# Vault Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a HashiCorp Vault client package. It uses the `hashicorp/vault/api` library and follows the standard module principles in `guide/module.md`.

### Architectural Role

*   **Location**: `pkg/client/vault`
*   **Purpose**: To provide a concrete implementation of a `lib` interface for secrets management, encapsulating all interactions with a Vault server.

### Generation Prompt for LLM

When asked to create a Vault client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/client/vault/fx.go`**

```go
package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
	"go.uber.org/fx"

	// e.g., "your/project/lib/secrets"
)

var Module = fx.Module("vault-client",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// fx.As(new(secrets.Manager)),
		),
		ConfigRegister,
	),
)

type Config struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}

func ConfigRegister() *Config {
	return &Config{
		Address: "http://127.0.0.1:8200",
	}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

type Client struct {
	*api.Client
}

func NewClient(p Param) (*Client, error) {
	config := &api.Config{
		Address: p.Config.Address,
	}

	vc, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	vc.SetToken(p.Config.Token)
	client := &Client{Client: vc}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := client.Sys().HealthStatus()
			return err
		},
		OnStop: nil,
	})

	return client, nil
}
```
