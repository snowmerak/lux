---
description: Guide for generating a gRPC server Fx module.
tags: [go, grpc, server, fx, module, controller, rpc]
---

# gRPC Server Fx Module Generation Guide

This guide is for generating a standardized Fx module for a gRPC server. This module typically resides in `internal/controller` and is responsible for exposing services via gRPC.

**File to Create**: `internal/controller/usergrpc/fx.go`

**LLM Prompt**:
"Create an Fx module for a gRPC server. The package name is `usergrpc`. It should depend on `user.Service` and `logger.Logger`. The module must register the gRPC server implementation and manage its lifecycle using `fx.Lifecycle`."

---

### Generated `fx.go`

```go
package usergrpc

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"your/project/lib/adapter/logger"
	pb "your/project/lib/adapter/gen/user" // Generated protobuf code
	"your/project/lib/service/user"
)

// Module exports the gRPC server component to the Fx application.
// We use fx.Invoke to register a function that will be executed on application start.
var Module = fx.Module("grpc-server",
	fx.Invoke(RegisterServer),
	ConfigRegister,
)

// Config holds the configuration for the gRPC server.
type Config struct {
	Port int `yaml:"port"`
}

// ConfigRegister provides the configuration for this module.
func ConfigRegister() *Config {
	return &Config{Port: 50051} // Default gRPC port
}

// Param defines the dependencies for the server.
type Param struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Logger      logger.Logger
	Config      *Config
	UserService user.Service // Depends on the user service
}

// server is the concrete implementation of the generated gRPC server interface.
type server struct {
	pb.UnimplementedUserServiceServer // Embed the unimplemented server
	userService user.Service
	logger      logger.Logger
}

// RegisterServer sets up and manages the gRPC server's lifecycle.
func RegisterServer(p Param) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", p.Config.Port))
	if err != nil {
		p.Logger.Error(err, "Failed to listen for gRPC server")
		return
	}

	grpcServer := grpc.NewServer()
	
	// Create an instance of our server implementation.
	srv := &server{
		userService: p.UserService,
		logger:      p.Logger,
	}
	// Register the server implementation with the gRPC server.
	pb.RegisterUserServiceServer(grpcServer, srv)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info("Starting gRPC server", map[string]any{"addr": lis.Addr().String()})
			// Run the server in a separate goroutine to avoid blocking.
			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					p.Logger.Error(err, "Failed to start gRPC server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Stopping gRPC server")
			grpcServer.GracefulStop()
			return nil
		},
	})
}

// Example gRPC method implementation.
// func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
// 	// ... implementation ...
// }
```
