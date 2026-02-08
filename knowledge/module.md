# Go Module Generation Guide (with Fx)

This document provides a comprehensive guide to creating standardized, modular components using the `uber-go/fx` framework, following the architectural principles outlined in `project_structure.md`. A consistent module structure is essential for building a scalable and maintainable application where dependencies are managed explicitly and lifecycles are handled gracefully.

### Core Principles of a Module

1.  **Self-Contained**: Each module should be a self-contained unit of functionality. It defines its own configuration, dependencies, and constructor.
2.  **Explicit Dependencies**: All dependencies required by a module **must** be declared explicitly using the `fx.In` struct tag within a `Param` struct. This makes the dependency graph clear and statically analyzable.
3.  **Interface-Driven**: Modules that provide implementations (typically in `pkg`) **must** provide them as `lib` interfaces using `fx.As`. This decouples the business logic from concrete implementations.
4.  **Lifecycle Aware**: Modules that manage resources (like database connections or message queue consumers) **must** use `fx.Lifecycle` to register `OnStart` and `OnStop` hooks for proper initialization and graceful shutdown.

### Module Templates by Layer

The application is divided into distinct layers (`pkg`, `internal/service`, `internal/controller`), and each has a specific role within the Fx dependency injection framework.

---

### 1. The `pkg` Module: The Provider

The `pkg` layer's role is to provide a concrete, technology-specific **implementation** of an interface defined in the `lib` layer. It acts as a "Provider" of functionality to the rest of the application.

**Key Characteristics:**
*   **Implements a `lib` interface.**
*   Uses `fx.Provide` with `fx.Annotate` and `fx.As` to bind the concrete type to its interface.
*   Manages the lifecycle of the resource it provides (e.g., a database connection).
*   **Must not** depend on other `lib` interfaces.

**Standard `fx.go` Template for a `pkg` Module:**

```go
// In pkg/client/postgres/fx.go
package postgres

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/fx"
	// Import the LIB interface that this module implements.
	"your/project/lib/repository/user"
)

// The compile-time interface check is mandatory.
var _ user.Repository = (*UserRepository)(nil)

// Module exports the component's functionality to the Fx application.
// The module name (e.g., "postgres") should be descriptive and unique within the application.
var Module = fx.Module("postgres-user-repo",
	// fx.Provide lists all the constructors this module offers to the DI container.
	fx.Provide(
		// The main constructor for the component.
		// fx.Annotate is used to explicitly associate the concrete implementation
		// with the interface it implements. This is a mandatory pattern.
		fx.Annotate(
			NewUserRepository, // The constructor function.
			// fx.As casts the concrete return type (*UserRepository) to one or more
			// interfaces (e.g., new(user.Repository)).
			fx.As(new(user.Repository)),
		),
		// The constructor for the module's configuration.
		ConfigRegister,
	),
)

// Config holds the configuration specific to this module.
// These values MUST be populated from environment variables.
type Config struct {
	DSN      string
	PoolSize int
}

// ConfigRegister loads and provides the module's configuration.
// It MUST enforce that required environment variables are present.
func ConfigRegister() (*Config, error) {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("POSTGRES_DSN environment variable is required")
	}

	poolSizeStr := os.Getenv("POSTGRES_POOL_SIZE")
	if poolSizeStr == "" {
		return nil, fmt.Errorf("POSTGRES_POOL_SIZE environment variable is required")
	}

	poolSize, err := strconv.Atoi(poolSizeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_POOL_SIZE: %w", err)
	}

	return &Config{
		DSN:      dsn,
		PoolSize: poolSize,
	}, nil
}

// Param is a struct that groups all dependencies for the main constructor.
// The `fx.In` tag tells Fx to populate the fields of this struct.
// This keeps the constructor signature clean and manageable.
type Param struct {
	fx.In

	// Lifecycle is a mandatory dependency for any module that manages a resource.
	Lifecycle fx.Lifecycle
	// The module's specific configuration.
	Config *Config
	// Other dependencies, such as a logger, can be added here.
	// Logger client.Logger
}

// UserRepository is the concrete struct that implements the user.Repository interface.
type UserRepository struct {
	// e.g., db *pgxpool.Pool
}

// NewUserRepository is the constructor for the UserRepository.
// It receives all its dependencies via the Param struct, provided by Fx.
func NewUserRepository(p Param) (*UserRepository, error) {
	repo := &UserRepository{
		// ... initialize fields ...
	}

	// The Fx lifecycle is used to register startup and shutdown hooks.
	// This is mandatory for any resource that needs to be initialized or cleaned up.
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Logic to run on application start.
			// e.g., connect to the database using p.Config.DSN
			// and assign the connection to repo.db.
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Logic to run on application stop.
			// e.g., close the database connection.
			return nil
		},
	})

	return repo, nil
}
```

---

### 2. The `internal/service` Module: The Composer

The `internal/service` layer is where core business logic lives. Its role is to **compose** multiple `lib` interfaces from different `pkg` providers to orchestrate complex business workflows.

**Key Characteristics:**
*   **Depends only on the `lib` package (interfaces and data structures).**
*   Does not depend on `internal/controller` or `pkg`.
*   Uses `fx.Provide` to make the service available to the `controller` layer.
*   Typically does not manage resources directly, so `fx.Lifecycle` is less common here.

**Standard `fx.go` Template for a `service` Module:**

```go
// In internal/service/membership/fx.go
package membership

import (
	"go.uber.org/fx"
	// Import necessary lib interfaces.
	"your/project/lib/adapter/notification"
	"your/project/lib/logger"
	"your/project/lib/repository/user"
)

var Module = fx.Module("membership-service",
	fx.Provide(
		NewService,
	),
)

// Param struct for the service. It depends ONLY on interfaces from `lib`.
type Param struct {
	fx.In

	UserRepo user.Repository
	Notifier notification.Adapter
	Logger   logger.Logger
}

// Service struct holds its dependencies, which are lib interfaces.
type Service struct {
	userRepo user.Repository
	notifier notification.Adapter
	logger   logger.Logger
}

// NewService is the constructor for the membership service.
func NewService(p Param) *Service {
	return &Service{
		userRepo: p.UserRepo,
		notifier: p.Notifier,
		logger:   p.Logger,
	}
}

// ... service methods that orchestrate business logic ...
// func (s *Service) RegisterUser(...) { ... }
```

---

### 3. The `internal/controller` Module: The Handler

The `internal/controller` layer acts as the entry point for external requests (e.g., HTTP, gRPC). Its primary role is to handle incoming data, call the appropriate `service` methods, and formulate a response. It does not contain business logic.

**Key Characteristics:**
*   **Depends on `internal/service` structs and the `lib` package.**
*   Uses `fx.Invoke` to register handlers (e.g., HTTP routes). Registration is a side effect, and the controller itself is not usually a dependency for other components, so `fx.Invoke` is preferred over `fx.Provide`.
*   The invoked function receives all dependencies needed to set up the handlers.

**Standard `fx.go` Template for a `controller` Module:**

```go
// In internal/controller/userapi/fx.go
package userapi

import (
	"net/http"

	"github.com/go-chi/chi/v5" // Example using chi router
	"go.uber.org/fx"

	"your/project/lib/logger"
	// Import the service it depends on.
	"your/project/internal/service/membership"
)

// Module uses fx.Invoke because a controller's primary role is to register
// handlers, which is a side effect, not providing a new dependency.
var Module = fx.Module("userapi-controller",
	fx.Invoke(RegisterRoutes),
)

// Param struct for the controller.
// It depends on services from `internal/service` and components via `lib` interfaces.
type Param struct {
	fx.In

	Router        chi.Router // Assumes a chi.Router is provided by a pkg module.
	MembershipSvc *membership.Service
	Logger        logger.Logger
}

// Controller holds dependencies needed by its handler methods.
type Controller struct {
	service *membership.Service
	logger  logger.Logger
}

// RegisterRoutes creates the controller and sets up its HTTP routes.
func RegisterRoutes(p Param) {
	c := &Controller{
		service: p.MembershipSvc,
		logger:  p.Logger,
	}
	p.Router.Post("/users", c.CreateUser)
	// ... other routes
}

// CreateUser is an example handler method.
func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Decode request.
	// 2. Call c.service to perform business logic.
	// 3. Encode and write response.
}
```

---

### Technology-Specific Examples

In addition to these standard templates, you can find detailed generation guides for specific technology stacks in the `guide/packages/` subdirectory. These guides provide concrete, ready-to-use Fx module examples for common components like databases (Postgres, Redis), message queues (Kafka, NATS), and HTTP servers (Chi, Fiber).

When writing a new module, it is recommended to first check the `guide/packages/` directory for a guide on a similar technology and use it as a baseline.

### Module Composition in `main.go`

As described in `entry_point.md`, the `cmd/.../app.go` file is responsible for assembling all the application's modules. Fx will build the dependency graph, execute the constructors in the correct order, and run the application.

```go
// In cmd/publicapi/app.go
package main

import (
	"go.uber.org/fx"

	// Import all necessary modules
	"your/project/internal/controller/userapi"
	"your/project/internal/service/membership"
	"your/project/pkg/client/chi"
	"your/project/pkg/client/postgres"
	"your/project/pkg/client/zerolog"
)

func main() {
	fx.New(
		// List all modules here. Fx resolves the dependency order.
		// pkg modules (Providers)
		chi.Module,
		postgres.Module,
		zerolog.Module,
		// service modules (Composers)
		membership.Module,
		// controller modules (Handlers)
		userapi.Module,
	).Run()
}
```

---

### Tooling Dependencies

Some `pkg` modules may require external command-line interface (CLI) tools for code generation (e.g., `sqlc`, `templ`). These tools are development-time dependencies and should not be part of the application's runtime.

**Installation**

These tools **must** be installed using the `go install` command to place the executable in your `$GOPATH/bin` directory. For example:

```sh
# Example for installing sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Example for installing templ
go install github.com/a-h/templ/cmd/templ@latest
```

**Documentation**

If a `pkg` module requires such a tool, its installation command and usage **must** be clearly documented in its corresponding `guide/packages/*.md` file. This ensures that other developers can easily set up their environment to work with the module.
