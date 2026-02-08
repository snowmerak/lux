# TimingWheel Package (`pkg`) Generation Guide

This guide provides a blueprint for creating a `timingwheel` module within the `pkg/` directory. It uses the `github.com/RussellLuo/timingwheel` library, which implements Hierarchical Timing Wheels for efficient timer management.

### Architectural Role

*   **Location**: `pkg/scheduler/timingwheel`
*   **Purpose**: To provide a high-performance scheduling mechanism for delayed tasks, cron jobs, or timeouts. This module encapsulates the timing wheel's lifecycle and exposes it via a `lib` interface.

### Generation Prompt for LLM

When asked to create a TimingWheel package, follow these steps.

#### 1. Create the Directory Structure

Create the directory `pkg/scheduler/timingwheel/`.

#### 2. Create the `fx.go` Module Definition

This file wires the TimingWheel into the application's dependency injection system.

**File: `pkg/scheduler/timingwheel/fx.go`**

```go
package timingwheel

import (
	"github.com/RussellLuo/timingwheel"
	"go.uber.org/fx"
)

// Module exports the TimingWheel functionality to the Fx application.
var Module = fx.Module("timingwheel-scheduler",
	fx.Provide(
		NewClient,
		ConfigRegister,
	),
)
```

#### 3. Create the `config.go` Configuration

**File: `pkg/scheduler/timingwheel/config.go`**

```go
package timingwheel

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Tick      time.Duration
	WheelSize int64
}

func ConfigRegister() (*Config, error) {
	tickStr := os.Getenv("TIMINGWHEEL_TICK_MS")
	if tickStr == "" {
		tickStr = "10" // Default 10ms
	}
	tickMs, err := strconv.ParseInt(tickStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TIMINGWHEEL_TICK_MS: %w", err)
	}

	wheelSizeStr := os.Getenv("TIMINGWHEEL_SIZE")
	if wheelSizeStr == "" {
		wheelSizeStr = "60" // Default 60 slots
	}
	wheelSize, err := strconv.ParseInt(wheelSizeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TIMINGWHEEL_SIZE: %w", err)
	}

	return &Config{
		Tick:      time.Duration(tickMs) * time.Millisecond,
		WheelSize: wheelSize,
	}, nil
}
```

#### 4. Create the `client.go` Implementation

This file manages the TimingWheel instance.

**File: `pkg/scheduler/timingwheel/client.go`**

```go
package timingwheel

import (
	"context"

	"github.com/RussellLuo/timingwheel"
	"go.uber.org/fx"
)

type Client struct {
	tw *timingwheel.TimingWheel
}

type Param struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *Config
}

func NewClient(p Param) *Client {
	tw := timingwheel.NewTimingWheel(p.Config.Tick, p.Config.WheelSize)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			tw.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			tw.Stop()
			return nil
		},
	})

	return &Client{tw: tw}
}

// After schedules a function to run after the specified duration.
func (c *Client) After(d time.Duration, fn func()) *timingwheel.Timer {
	return c.tw.After(d, fn)
}

// ScheduleFunc schedules a function to run based on a scheduler.
func (c *Client) ScheduleFunc(s timingwheel.Scheduler, f func()) *timingwheel.Timer {
	return c.tw.ScheduleFunc(s, f)
}
```

#### 5. Implement a Scheduler (Example)

If implementing a specific scheduling interface from `lib/adapter/scheduler.go`:

**File: `pkg/scheduler/timingwheel/adapter.go`**

```go
package timingwheel

import (
	"time"
	"your/project/lib/adapter"
)

var _ adapter.Scheduler = (*schedulerAdapter)(nil)

type schedulerAdapter struct {
	client *Client
}

func NewSchedulerAdapter(client *Client) adapter.Scheduler {
	return &schedulerAdapter{client: client}
}

func (a *schedulerAdapter) Schedule(d time.Duration, task func()) {
	a.client.After(d, task)
}
```
