# MinIO Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a MinIO client package. It uses the `minio/minio-go/v7` library and follows the standard module principles in `guide/module.md`.

### Architectural Role

*   **Location**: `pkg/client/minio`
*   **Purpose**: To provide a concrete implementation of a `lib` interface for S3-compatible object storage.

### Generation Prompt for LLM

When asked to create a MinIO client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/client/minio/fx.go`**

```go
package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"

	// e.g., "your/project/lib/storage"
)

var Module = fx.Module("minio-client",
	fx.Provide(
		fx.Annotate(
			NewClient,
			// fx.As(new(storage.ObjectStorage)),
		),
		ConfigRegister,
	),
)

type Config struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `yaml:"useSsl"`
}

func ConfigRegister() *Config {
	return &Config{
		Endpoint: "localhost:9000",
	}
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

type Client struct {
	*minio.Client
}

func NewClient(p Param) (*Client, error) {
	mc, err := minio.New(p.Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(p.Config.AccessKeyID, p.Config.SecretAccessKey, ""),
		Secure: p.Config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	client := &Client{Client: mc}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := client.ListBuckets(ctx)
			return err
		},
		OnStop: nil,
	})

	return client, nil
}
```
