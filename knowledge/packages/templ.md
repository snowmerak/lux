---
description: Guide for generating a Templ renderer Fx module for server-side rendering.
tags: [go, templ, html, renderer, fx, module]
---

# Templ Renderer Fx Module Generation Guide

This guide is for generating a standardized Fx module for a `templ` renderer. This allows `templ` components to be used for server-side rendering within an HTTP controller (like Chi or Fiber).

**File to Create**: `pkg/renderer/templ/fx.go`

**LLM Prompt**:
"Create an Fx module for a `templ` renderer. The package name is `templ`. It should implement a `lib/renderer.Renderer` interface. This module will allow `templ` components to be rendered in HTTP handlers."

---

### `lib/renderer/renderer.go` (Interface Definition)

First, define a generic interface for a renderer in the `lib` layer.

```go
package renderer

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// Renderer defines a standard interface for rendering components.
type Renderer interface {
	Render(ctx context.Context, w io.Writer, component templ.Component) error
}
```

### Generated `fx.go`

```go
package templ

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"go.uber.org/fx"

	"your/project/lib/renderer"
)

// Module exports the templ renderer component to the Fx application.
var Module = fx.Module("templ-renderer",
	fx.Provide(
		fx.Annotate(
			NewRenderer,
			fx.As(new(renderer.Renderer)),
		),
	),
)

// templRenderer is the concrete implementation of the renderer.Renderer interface.
type templRenderer struct{}

// NewRenderer creates a new templ renderer.
// This component is stateless, so it has no dependencies or config.
func NewRenderer() (renderer.Renderer, error) {
	// Compile-time interface check.
	var _ renderer.Renderer = (*templRenderer)(nil)
	return &templRenderer{}, nil
}

// Render executes the given templ component and writes the output.
func (r *templRenderer) Render(ctx context.Context, w io.Writer, component templ.Component) error {
	return component.Render(ctx, w)
}
```

### Usage in an HTTP Handler (Example)

```go
// In your internal/controller/userapi/controller.go

// ...
import (
	"your/project/lib/renderer"
	// your generated templ components
	// "your/project/internal/controller/userapi/view"
)

// ...

func makeGetUserPageHandler(renderer renderer.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ... fetch user data
		// component := view.Profile(user)
		// renderer.Render(r.Context(), w, component)
	}
}
```

### Prerequisites

1.  **Chi Module:** You need a `chi` module that provides a `chi.Router`. Refer to the `chi.md` guide for creating one.

### Tooling

`templ` is a CLI tool that generates Go code from `.templ` files. You must have it available to generate the necessary Go code for your components.

#### Tool Dependency Management (Go 1.24+ Recommended)

With Go 1.24 and later, you can manage tool dependencies directly within your `go.mod` file using the `go get -tool` command. This is the official and recommended way to ensure that all developers and CI environments use the exact same version of a tool.

1.  **Add the tool to `go.mod`:**
    Run the following command to add `templ` as a tool dependency. The `-tool` flag tells the `go` command to add it as a development tool, not a regular dependency.

    ```sh
    go get -tool github.com/a-h/templ/cmd/templ
    ```

    This will add a `tool` directive to your `go.mod` file, similar to this:
    ```
    toolchain go1.24.0

    tool github.com/a-h/templ/cmd/templ v0.2.648
    ```

2.  **Usage:**
    You can now run `templ` using the `go tool` command, which executes the specific version defined in your `go.mod`.

    ```sh
    go tool templ generate
    ```

### How It Works
