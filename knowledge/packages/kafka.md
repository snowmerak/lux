# Kafka Client Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a Kafka client package. It uses the `segmentio/kafka-go` library and follows the standard module principles in `guide/module.md`.

### Architectural Role

*   **Location**: `pkg/adapter/kafka`
*   **Purpose**: To provide a concrete implementation for producing and consuming messages from Kafka topics.

### Generation Prompt for LLM

When asked to create a Kafka client package, follow these steps.

#### 1. Create the `fx.go` Module Definition

**File: `pkg/adapter/kafka/fx.go`**

```go
package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

var Module = fx.Module("kafka-adapter",
	fx.Provide(
		NewWriter, // Provides a Kafka Writer (Producer)
		ConfigRegister,
	),
)

type Config struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

func ConfigRegister() *Config {
	return &Config{
		Brokers: []string{"localhost:9092"},
		Topic:   "default-topic",
	}
}

type WriterParam struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

// NewWriter creates a new Kafka writer (producer).
func NewWriter(p WriterParam) *kafka.Writer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(p.Config.Brokers...),
		Topic:    p.Config.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			return w.Close()
		},
	})

	return w
}

// Note: A Kafka Reader (Consumer) would typically be created and started
// inside an fx.Invoke call within a service or controller module that needs to consume messages.
```
