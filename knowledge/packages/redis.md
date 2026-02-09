---
description: Guide for generating a Redis client Fx module for caching and session management.
tags: [go, redis, cache, fx, module]
---

# Redis Client Package (`pkg`) Generation Guide

This guide provides a detailed blueprint for creating a Redis client package within the `pkg/client/` directory. It follows the standard module principles defined in `guide/module.md` and is specifically tailored for implementing a Redis client using the `redis/go-redis` library.

### Architectural Role

*   **Location**: `pkg/client/redis`
*   **Purpose**: To provide a concrete implementation of one or more `lib` interfaces that require Redis for caching, session management, etc. This module encapsulates all interactions with the Redis server.

### Generation Prompt for LLM

When asked to create a Redis client package, follow these steps to generate the required files.

#### 1. Create the Directory Structure

Create the directory `pkg/client/redis/`.

#### 2. Create the `fx.go` Module Definition

This file wires the Redis client into the application's dependency injection system. It provides the client constructor and its configuration.

**File: `pkg/client/redis/fx.go`**

```go
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	// Import the specific lib interface(s) this client will implement.
	// For example:
	// "your/project/lib/cache"
)

// Module exports the Redis client functionality to the Fx application.
var Module = fx.Module("redis-client",
	fx.Provide(
		// Provide the constructor for the Redis client.
		// It is annotated to be cast to the specific interface it implements.
		fx.Annotate(
			NewClient,
			// fx.As(new(cache.Cache)), // Example: cast to a cache interface
		),
		// Provide the configuration constructor.
		ConfigRegister,
	),
)

// Config holds the configuration for the Redis client.
type Config struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// ConfigRegister loads the Redis configuration.
func ConfigRegister() *Config {
	// In a real application, load from a file or environment variables.
	return &Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

// Param groups the dependencies for the Redis client constructor.
type Param struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
}

// Client is the concrete implementation for Redis.
type Client struct {
	*redis.Client
}

// NewClient is the constructor for the Redis client.
func NewClient(p Param) (*Client, error) {
	// The compile-time interface check is mandatory.
	// var _ cache.Cache = (*Client)(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr:     p.Config.Addr,
		Password: p.Config.Password,
		DB:       p.Config.DB,
	})

	client := &Client{Client: rdb}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Ping the Redis server on startup to ensure the connection is valid.
			_, err := rdb.Ping(ctx).Result()
			return err
		},
		OnStop: func(ctx context.Context) error {
			// Gracefully close the Redis connection.
			return rdb.Close()
		},
	})

	return client, nil
}
```

#### 3. Implement the Interface Methods

If the `Client` struct implements a `lib` interface (as it should), the methods for that interface should be implemented in a separate file, for example `cache.go`.

**File: `pkg/client/redis/cache.go`**

```go
package redis

import "context"

// Example implementation of a Set method for a cache.Cache interface.
/*
func (c *Client) Set(ctx context.Context, key string, value interface{}) error {
    return c.Client.Set(ctx, key, value, 0).Err()
}
*/
```
