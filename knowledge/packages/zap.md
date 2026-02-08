# Zap Logger Generation Guide

This guide provides a template for creating a logger module using `uber-go/zap`, which is a popular and highly performant logging library. This module will be integrated into the application using Fx.

This implementation is an alternative to the `zerolog` implementation detailed in the main `logger.md` guide.

## Module Definition

```go
package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"{{.gomod}}/lib/logger"
)

var Module = fx.Module("zap-logger",
	fx.Provide(
		NewZapLogger,
	),
)

// NewZapLogger creates a new zap.Logger instance.
// It can be configured to read from a config file or environment variables.
func NewZapLogger() (*zap.Logger, error) {
	// For development, a development logger is often useful.
	// For production, you would likely use zap.NewProduction() or a custom configuration.
	return zap.NewDevelopment()
}

// Provide this in your main logger module to satisfy the logger.Logger interface
func NewLogger(log *zap.Logger) logger.Logger {
    return &ZapLogger{logger: log}
}

// ZapLogger is an adapter to satisfy the logger.Logger interface.
type ZapLogger struct {
	logger *zap.Logger
}

func (l *ZapLogger) Info(msg string, args ...interface{}) {
	l.logger.Sugar().Info(msg, args...)
}

func (l *ZapLogger) Error(msg string, args ...interface{}) {
	l.logger.Sugar().Error(msg, args...)
}

func (l *ZapLogger) Debug(msg string, args ...interface{}) {
	l.logger.Sugar().Debug(msg, args...)
}

func (l *ZapLogger) Warn(msg string, args ...interface{}) {
	l.logger.Sugar().Warn(msg, args...)
}

func (l *ZapLogger) GetZapLogger() *zap.Logger {
    return l.logger
}

func (l *ZapLogger) GetZerologLogger() *zerolog.Logger {
    return nil // Not applicable
}
```

### Integration

To use this `zap` logger, you would include its `Module` in your Fx application. You would also need to provide the `NewLogger` function in your main logger module to bind the `*zap.Logger` to the `logger.Logger` interface.
