# SQLC Repository Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a repository package using `sqlc`. It follows the standard module principles in `guide/module.md` and the project structure in `guide/project_structure.md`.

### Architectural Role

*   **Location**: e.g., `pkg/repository/user`
*   **Purpose**: To provide a concrete, type-safe implementation of a `lib/repository` interface. It uses `sqlc` to generate Go code from raw SQL queries and interacts with a database via a standard library driver (e.g., `pgx`).

### Generation Prompt for LLM

When asked to create a repository implementation with `sqlc`, follow these steps precisely.

#### 1. Create the Directory Structure

Create the following directories inside the package (e.g., `pkg/repository/user`):

*   `sql/`: To hold all SQL-related files.
*   `sql/migrations/`: For database schema definitions.
*   `sql/queries/`: For `sqlc` queries.

#### 2. Create `sqlc.yaml` Configuration

Place this file in the root of the package directory (e.g., `pkg/repository/user/sqlc.yaml`).

```yaml
version: "2"
sql:
  - engine: "postgresql" # Or "mysql"
    queries: "sql/queries/"
    schema: "sql/migrations/"
    gen:
      go:
        package: "user" # Should match the parent directory name
        out: "." # Generate code in the current directory
        sql_package: "pgx/v5"
        emit_interface: true # Recommended: generates an interface for the querier
```

#### 3. Create SQL Files

*   **Schema**: `sql/migrations/001_users.sql`
    ```sql
    CREATE TABLE users (
        id UUID PRIMARY KEY,
        name TEXT NOT NULL
    );
    ```
*   **Queries**: `sql/queries/user.sql`
    ```sql
    -- name: GetUser :one
    SELECT * FROM users
    WHERE id = $1 LIMIT 1;

    -- name: CreateUser :one
    INSERT INTO users (id, name)
    VALUES ($1, $2)
    RETURNING *;
    ```

#### 4. Run `sqlc` to Generate Code

Execute `sqlc generate` within the package directory. This will create `db.go`, `models.go`, and `user.sql.go` based on the configuration.

#### 5. Create the `fx.go` Module Definition

This file provides the repository implementation to the Fx application. It depends on a database connection pool, which is typically provided by another module (e.g., a `pkg/client/postgres` module).

**File: `pkg/repository/user/fx.go`**

```go
package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"your/project/lib/repository/user" // Import the LIB interface
)

// Module exports the repository to the Fx application.
var Module = fx.Module("user-repository",
	fx.Provide(
		fx.Annotate(
			NewRepository,
			fx.As(new(user.Repository)), // Cast to the lib interface
		),
	),
)

// Param groups the dependencies for the repository constructor.
type Param struct {
	fx.In

	// This module depends on a database connection pool.
	Pool *pgxpool.Pool
}

// Repository is the concrete implementation of the user.Repository interface.
type Repository struct {
	*Queries // Embed the generated sqlc Querier
	pool     *pgxpool.Pool
}

// NewRepository is the constructor for the repository.
func NewRepository(p Param) (user.Repository, error) {
	// The compile-time interface check is mandatory.
	var _ user.Repository = (*Repository)(nil)

	return &Repository{
		Queries: New(p.Pool),
		pool:    p.Pool,
	}, nil
}
```

#### 6. Implement the Interface

Create a file like `repository.go` to implement the methods defined in your `lib` interface, using the generated `sqlc` methods.

**File: `pkg/repository/user/repository.go`**

```go
package user

import (
	"context"
	"github.com/google/uuid"
	"your/project/lib/domain"
)

// GetUser implements the user.Repository interface.
func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	// Call the generated sqlc method.
	user, err := r.Queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// Map the generated model to the domain model.
	return &domain.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

// ... other interface methods
```

### Tooling

`sqlc` is a code generator that creates type-safe Go code from SQL. You must install it to use this module.

#### Tool Dependency Management (Go 1.24+ Recommended)

With Go 1.24 and later, you can manage tool dependencies directly within your `go.mod` file using the `go get -tool` command. This is the official and recommended way to ensure that all developers and CI environments use the exact same version of a tool.

1.  **Add the tool to `go.mod`:**
    Run the following command to add `sqlc` as a tool dependency. The `-tool` flag tells the `go` command to add it as a development tool, not a regular dependency.

    ```sh
    go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc
    ```

    This will add a `tool` directive to your `go.mod` file, similar to this:
    ```
    toolchain go1.24.0

    tool github.com/sqlc-dev/sqlc/cmd/sqlc v1.26.0
    ```

2.  **Usage:**
    You can now run `sqlc` using the `go tool` command, which executes the specific version defined in your `go.mod`.

    ```sh
    go tool sqlc generate
    ```

    You should run this command from the directory containing your `sqlc.yaml` file.
