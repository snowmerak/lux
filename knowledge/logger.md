# Logging Guide

This document defines two primary patterns for implementing consistent, structured logging across the application. All logging must be performed through a standard logger interface injected via `uber-go/fx`.

### Core Principles

1.  **Structured Logging**: All logs must be structured as `key-value` pairs. This allows for easy parsing, searching, and analysis.
2.  **Interface-Driven**: All business logic (`internal/service`) must depend on the `Logger` interface defined in `lib/adapter/logger`, not on a concrete implementation.
3.  **Context Propagation**: When possible, propagate information like request-specific IDs or trace IDs via `context.Context` and include it in logs.

---

### 1. Logger Interface and Implementation

#### `lib/adapter/logger/logger.go`

This is the standard logger interface that all components of the application must depend on.

```go
package logger

import "context"

// Logger defines the standard logging interface for the application.
// It is designed to be a simple, structured logging API.
type Logger interface {
	Debug(msg string, fields ...map[string]any)
	Info(msg string, fields ...map[string]any)
	Warn(msg string, fields ...map[string]any)
	Error(err error, msg string, fields ...map[string]any)
	// With returns a new logger with structured context.
	With(fields map[string]any) Logger
	// Ctx returns a logger that extracts trace information from the context.
	Ctx(ctx context.Context) Logger
}
```

#### `pkg/adapter/zerolog/fx.go`

This is the concrete implementation of the `Logger` interface using `zerolog` and its Fx module.

```go
package zerolog

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"your/project/lib/adapter/logger"
)

// Module provides the zerolog logger to the Fx application.
var Module = fx.Module("zerolog",
	fx.Provide(
		// Provide NewZerologLogger as an implementation of logger.Logger.
		fx.Annotate(
			NewZerologLogger,
			fx.As(new(logger.Logger)),
		),
	),
)

// zerologLogger is the zerolog implementation of the logger.Logger interface.
type zerologLogger struct {
	logger zerolog.Logger
}

// NewZerologLogger creates a new zerolog logger.
func NewZerologLogger() logger.Logger {
	// You can replace this with a different writer, e.g., for JSON output.
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &zerologLogger{logger: logger}
}

func (l *zerologLogger) log(e *zerolog.Event, msg string, fields ...map[string]any) {
	if len(fields) > 0 && fields[0] != nil {
		e.Fields(fields[0]).Msg(msg)
	} else {
		e.Msg(msg)
	}
}

func (l *zerologLogger) Debug(msg string, fields ...map[string]any) {
	l.log(l.logger.Debug(), msg, fields...)
}

func (l *zerologLogger) Info(msg string, fields ...map[string]any) {
	l.log(l.logger.Info(), msg, fields...)
}

func (l *zerologLogger) Warn(msg string, fields ...map[string]any) {
	l.log(l.logger.Warn(), msg, fields...)
}

func (l *zerologLogger) Error(err error, msg string, fields ...map[string]any) {
	event := l.logger.Error().Err(err)
	l.log(event, msg, fields...)
}

func (l *zerologLogger) With(fields map[string]any) logger.Logger {
	return &zerologLogger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}

func (l *zerologLogger) Ctx(ctx context.Context) logger.Logger {
	return &zerologLogger{
		logger: zerolog.Ctx(ctx).With().Logger(),
	}
}
```

---

### 2. Logging Patterns

#### Pattern 1: Execution Trace Logging

Record the start and end of a method to trace the execution flow, duration, and success/failure status. Using a `defer` statement ensures that the exit log is always recorded, keeping the code clean.

**`internal/service/user/service.go` Example:**

```go
package user

import (
	"context"
	"time"
	"your/project/lib/adapter/logger"
)

type Service struct {
	logger logger.Logger
	// ... other dependencies
}

func (s *Service) CreateUser(ctx context.Context, email string) (id string, err error) {
	// Create a logger with method-specific context.
	log := s.logger.Ctx(ctx).With(map[string]any{
		"method": "CreateUser",
		"email":  email,
	})

	log.Info("Executing method")
	startTime := time.Now()

	// Use defer to ensure this logic runs right before the method returns.
	defer func() {
		duration := time.Since(startTime)
		fields := map[string]any{"duration_ms": duration.Milliseconds()}
		
		// Log differently based on whether an error occurred.
		if err != nil {
			log.Error(err, "Method execution failed", fields)
		} else {
			fields["user_id"] = id
			log.Info("Method execution successful", fields)
		}
	}()

	// ... business logic ...
	// id, err = s.userRepository.Create(...)
	
	return "new-user-id", nil
}
```

#### Pattern 2: State Change Auditing

When the state of important data (an entity) changes during method execution, log the before and after values to create an audit trail.

**`internal/service/user/service.go` Example:**

```go
func (s *Service) UpdateUserName(ctx context.Context, userID, newName string) (err error) {
	log := s.logger.Ctx(ctx).With(map[string]any{
		"method":  "UpdateUserName",
		"user_id": userID,
	})
	
	// ...

	// 1. Fetch the previous state.
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	oldName := user.Name

	// 2. Change the state.
	user.SetName(newName)

	// 3. Log the state change.
	log.Info("User name changed", map[string]any{
		"change_type": "Update",
		"field":       "Name",
		"before":      oldName,
		"after":       newName,
	})

	// 4. Save the new state.
	err = s.userRepository.Update(ctx, user)
	return err
}
```

### Fx Dependency Injection

Explicitly request `logger.Logger` in the `internal/service` constructor to receive it from Fx.

**`internal/service/user/fx.go` Example:**

```go
// ...
import "your/project/lib/adapter/logger"

type Param struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    logger.Logger // Request the logger dependency
	// ...
}

func NewService(p Param) (*Service, error) {
	service := &Service{
		logger: p.Logger,
		// ...
	}
	// ...
	return service, nil
}
```
