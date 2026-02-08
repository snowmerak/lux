# Fiber HTTP Server Fx Module Generation Guide

This guide is for generating a standardized Fx module for a Fiber-based HTTP server using Fiber v3. This module typically resides in `internal/controller` and is responsible for routing requests to services.

**File to Create**: `internal/controller/userapi/fx.go`

**LLM Prompt**:
"Create an Fx module for a Fiber v3 HTTP server. The package name is `userapi`. It should depend on `user.Service` and `logger.Logger`. The module must register the server routes and manage the HTTP server lifecycle using `fx.Lifecycle`."

---

### Generated `fx.go`

```go
package userapi

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"

	"your/project/lib/adapter/logger"
	"your/project/lib/service/user"
)

// Module exports the Fiber HTTP server component to the Fx application.
// We use fx.Invoke to register a function that will be executed on application start.
var Module = fx.Module("fiber-server",
	fx.Invoke(RegisterServer),
	ConfigRegister,
)

// Config holds the configuration for the HTTP server.
type Config struct {
	Port int `yaml:"port"`
}

// ConfigRegister provides the configuration for this module.
func ConfigRegister() *Config {
	return &Config{Port: 8080} // Default port
}

// Param defines the dependencies for the server.
type Param struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Logger      logger.Logger
	Config      *Config
	UserService user.Service // Depends on the user service
}

// RegisterServer sets up the Fiber app and manages the HTTP server's lifecycle.
func RegisterServer(p Param) {
	app := fiber.New()

	// Register routes
	// In Fiber v3, fiber.Ctx is an interface, so use 'c fiber.Ctx' instead of 'c *fiber.Ctx'.
	// Example: app.Get("/users/:id", func(c fiber.Ctx) error { ... })

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", p.Config.Port)
			p.Logger.Info("Starting HTTP server", map[string]any{"addr": addr})
			// Run the server in a separate goroutine to avoid blocking.
			go func() {
				if err := app.Listen(addr); err != nil {
					p.Logger.Error(err, "Failed to start HTTP server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Stopping HTTP server")
			return app.Shutdown()
		},
	})
}
```

### Key Changes in Fiber v3

1.  **Import Path**: Changed from `github.com/gofiber/fiber/v2` to `github.com/gofiber/fiber/v3`.
2.  **`fiber.Ctx` Interface**: `fiber.Ctx` is now an interface, not a pointer to a struct. Use `c fiber.Ctx`.
3.  **Binding API**: Unified binding API via `c.Bind()`.
    *   `c.BodyParser()` -> `c.Bind().Body()`
    *   `c.QueryParser()` -> `c.Bind().Query()`
    *   `c.ParamsParser()` -> `c.Bind().URI()`
4.  **Context Integration**: `fiber.Ctx` now implements `context.Context` directly.
5.  **Hooks**: Enhanced lifecycle hooks are available via `app.Hooks()`.
