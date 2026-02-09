---
description: Guide for generating a Chi HTTP server Fx module.
tags: [go, chi, http, server, fx, module, controller]
---

# Chi HTTP Server Fx Module Generation Guide

This guide is for generating a standardized Fx module for a Chi-based HTTP server. This module typically resides in `internal/controller` and is responsible for routing requests to services.

**File to Create**: `internal/controller/userapi/fx.go`

**LLM Prompt**:
"Create an Fx module for a Chi HTTP server. The package name is `userapi`. It should depend on `user.Service` and `logger.Logger`. The module must register the server routes and manage the HTTP server lifecycle using `fx.Lifecycle`."

---

### Generated `fx.go`

```go
package userapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"your/project/lib/adapter/logger"
	"your/project/lib/service/user"
)

// Module exports the Chi HTTP server component to the Fx application.
// Instead of fx.Provide, we use fx.Invoke to register a function that
// will be executed on application start. This is ideal for starting servers.
var Module = fx.Module("chi-server",
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

// RegisterServer sets up the Chi router and manages the HTTP server's lifecycle.
func RegisterServer(p Param) {
	router := chi.NewRouter()

	// Register routes
	// Example: router.Get("/users/{id}", makeGetUserHandler(p.UserService))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", p.Config.Port),
		Handler: router,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info("Starting HTTP server", map[string]any{"addr": server.Addr})
			// Run the server in a separate goroutine to avoid blocking.
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					p.Logger.Error(err, "Failed to start HTTP server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})
}
```
