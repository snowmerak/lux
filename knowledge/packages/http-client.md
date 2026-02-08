# HTTP Client Fx Module Generation Guide

This guide is for generating a standardized Fx module for a generic HTTP client. This module is useful for creating reusable clients in `pkg` that communicate with external REST APIs.

**File to Create**: `pkg/client/myapi/fx.go`

**LLM Prompt**:
"Create an Fx module for a generic HTTP client for 'MyAPI'. The package name is `myapi`. It should implement the `lib/adapter/http.Client` interface. The configuration should include a `BaseURL` and a `Timeout` duration. The module should be a standard Fx provider."

---

### `lib/adapter/http/client.go` (Interface Definition)

First, define a generic interface for the client in the `lib` layer.

```go
package http

import "net/http"

// Client defines a standard interface for an HTTP client.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}
```

### Generated `fx.go`

```go
package myapi

import (
	"net/http"
	"time"

	"go.uber.org/fx"
	libhttp "your/project/lib/adapter/http"
)

// Module exports the HTTP client component to the Fx application.
var Module = fx.Module("http-client-myapi",
	fx.Provide(
		fx.Annotate(
			NewClient,
			fx.As(new(libhttp.Client)),
		),
		ConfigRegister,
	),
)

// Config holds the configuration for the HTTP client.
type Config struct {
	BaseURL string        `yaml:"baseUrl"`
	Timeout time.Duration `yaml:"timeout"`
}

// ConfigRegister provides the configuration for this module.
func ConfigRegister() *Config {
	return &Config{
		Timeout: 30 * time.Second, // Default timeout
	}
}

// Param defines the dependencies for the client constructor.
type Param struct {
	fx.In

	Config *Config
	// Unlike servers, a simple client often doesn't need fx.Lifecycle,
	// as http.Client doesn't require explicit start/stop operations.
}

// client is the concrete implementation of the http.Client interface.
type client struct {
	*http.Client
	baseURL string
}

// NewClient creates a new HTTP client.
func NewClient(p Param) (libhttp.Client, error) {
	// Compile-time interface check.
	var _ libhttp.Client = (*client)(nil)

	return &client{
		Client: &http.Client{
			Timeout: p.Config.Timeout,
		},
		baseURL: p.Config.BaseURL,
	}, nil
}

// You can add methods that use the baseURL here.
// func (c *client) GetUser(id string) (*http.Response, error) {
// 	return c.Get(fmt.Sprintf("%s/users/%s", c.baseURL, id))
// }
```
