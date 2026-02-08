# Go CMD Entry Point Generation Guide (with Fx)

This guide provides a blueprint for generating the `main` package for a Go application, located within the `cmd` directory. The core of this architecture relies on the `uber-go/fx` framework to manage dependency injection and the application lifecycle. Following this guide will produce a clean, modular, and maintainable application entry point.

### Core Principles

1.  **Single Responsibility**: The `cmd` layer's **only** responsibility is to "wire" the application together. It connects concrete implementations from the `pkg` directory to the interfaces defined in `lib` and injects them into the business logic layers (`internal/service`, `internal/controller`). It should contain no business logic.
2.  **Fx as the Foundation**: `fx.App` is the heart of the application. The `main` function should be minimal, typically only containing the `fx.New(...).Run()` call. Fx handles the complexity of initializing components in the correct order and managing their lifecycles.
3.  **Modularity with `fx.Module`**: Every logical unit of the application (e.g., a database client, a repository, a service, a controller) **must** be provided as a separate `fx.Module`. This makes the dependency graph explicit, prevents circular dependencies, and allows for easy replacement of implementations.

### Generation Prompt for LLM

When asked to create a `cmd` entry point for a Go application, follow these steps precisely:

#### 1. Create the Application Directory

Create a new directory inside `cmd/` that represents the runnable application. For example:
*   `cmd/api` for a public REST API server.
*   `cmd/grpc` for a gRPC server.
*   `cmd/worker` for a background processing worker.

#### 2. Define a Standardized `fx.go` for Each Component

For each component in `pkg`, `internal/service`, and `internal/controller`, create a dedicated `fx.go` file that exposes its functionality as a standardized `fx.Module`.

**Standard `fx.go` Template:**

This template should be adapted for each component.

```go
// The package name should match the directory name.
package postgres 

import (
	"context"

	"go.uber.org/fx"
	"your/project/lib/repository/user" // Import the LIB interface
)

// Module exports the component's functionality to the Fx application.
// The module name should be descriptive and unique.
var Module = fx.Module("postgres",
	// Provide the constructor for the component.
	// fx.Annotate is used to explicitly declare that the concrete type (*UserRepository)
	// implements the interface (user.Repository). This is crucial for the DI graph.
	fx.Provide(
		fx.Annotate(
			NewUserRepository,
			fx.As(new(user.Repository)), // Cast to the interface
		),
		ConfigRegister, // Provide the configuration constructor
	),
)

// Config holds the configuration for this specific module.
type Config struct {
	// Configuration fields go here (e.g., DSN, PoolSize)
}

// ConfigRegister provides the config constructor to Fx.
// In a real application, this would load from a file or environment variables.
func ConfigRegister() *Config {
	return &Config{}
}

// Param groups the dependencies for the constructor using fx.In.
// This makes the constructor signature clean and easy to read.
type Param struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
	// Add other dependencies here, e.g., a logger client
}

// UserRepository is the concrete implementation of the user.Repository interface.
type UserRepository struct {
	// fields for the implementation
}

// NewUserRepository is the constructor for the UserRepository.
// It receives its dependencies via the Param struct.
func NewUserRepository(p Param) *UserRepository {
	repo := &UserRepository{}

	// The Fx lifecycle is used to register startup and shutdown hooks.
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Logic to run on application start, e.g., connect to the database.
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Logic to run on application stop, e.g., close database connections.
			return nil
		},
	})

	return repo
}
```

#### 3. Create the Server/Runner Module

Create a module responsible for the application's primary task, such as running an HTTP server. This module uses `fx.Invoke` to register a function that depends on other components (like controllers) and uses `fx.Lifecycle` to manage its own lifecycle.

**Example: `cmd/api/server/fx.go`**

```go
package server

import (
	"context"
	"net/http"
	"go.uber.org/fx"
	"your/project/internal/controller/userapi" // Controller dependency
)

var Module = fx.Module("server",
	fx.Provide(NewRouter),
	fx.Invoke(func(lc fx.Lifecycle, router *http.ServeMux, userController *userapi.Controller) {
		userController.RegisterRoutes(router) // Register routes from the controller

		server := &http.Server{ Addr: ":8080", Handler: router }

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go server.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		})
	}),
)

func NewRouter() *http.ServeMux {
	return http.NewServeMux()
}
```

#### 4. Assemble the Final `fx.App` in `main.go`

In the application directory (e.g., `cmd/api/`), create a `main.go` file that assembles all the modules and runs the `fx.App`.

**Example: `cmd/api/main.go`**

```go
package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"
	
	// Import all necessary modules
	"your/project/cmd/api/server"
	"your/project/internal/controller/userapi"
	"your/project/internal/service/user"
	"your/project/pkg/client/postgres"
	"your/project/pkg/adapter/logger"
)

func newStartupContext() context.Context {
	// A startup context with a timeout can prevent the app from hanging.
	return context.Background()
}

func main() {
	app := fx.New(
		// Provide a startup context.
		fx.Provide(newStartupucator),

		// Set timeouts for application start and stop.
		fx.StartTimeout(15*time.Second),
		fx.StopTimeout(15*time.Second),

		// Import all the modules that make up the application.
		// The order generally doesn't matter as Fx resolves the dependency graph.
		logger.Module,
		postgres.Module,
		user.Module,
		userapi.Module,
		server.Module,
	)

	// Run the application. This call blocks until the application is stopped.
	app.Run()

	// Check if the application exited with an error.
	if err := app.Err(); err != nil {
		log.Printf("application exited with error: %v", err)
	}
}
```
