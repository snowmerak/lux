# gRPC Client Fx Module Generation Guide

This guide is for generating a standardized Fx module for a gRPC client. The module will manage the gRPC connection lifecycle and provide a typed client as a dependency to other parts of the application.

**File to Create**: `pkg/client/servicea/fx.go`

**LLM Prompt**:
"Create an Fx module for a gRPC client connecting to 'Service A'. The package name is `servicea`. It should provide a `servicea.ServiceClient` (the generated gRPC client). The configuration should include the server `Address`. The module must manage the gRPC connection lifecycle using `fx.Lifecycle`."

---

### Generated `fx.go`

```go
package servicea

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import the generated protobuf client code
	pb "your/project/lib/adapter/gen/servicea"
)

// Module exports the gRPC client component to the Fx application.
var Module = fx.Module("grpc-client-servicea",
	fx.Provide(
		// This provides the concrete pb.ServiceClient implementation.
		// In a real app, you might wrap this in your own struct and provide
		// it as an interface using fx.As.
		NewServiceClient,
		ConfigRegister,
	),
)

// Config holds the configuration for the gRPC client.
type Config struct {
	Address string `yaml:"address"`
}

// ConfigRegister provides the configuration for this module.
func ConfigRegister() *Config {
	return &Config{}
}

// Param defines the dependencies for the gRPC client constructor.
type Param struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
}

// NewServiceClient creates a new gRPC service client and manages its lifecycle.
func NewServiceClient(p Param) (pb.ServiceClient, error) {
	var conn *grpc.ClientConn

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			// In production, you should use secure credentials.
			creds := grpc.WithTransportCredentials(insecure.NewCredentials())
			conn, err = grpc.DialContext(ctx, p.Config.Address, creds)
			return err
		},
		OnStop: func(ctx context.Context) error {
			if conn != nil {
				return conn.Close()
			}
			return nil
		},
	})

	// The connection is established in OnStart, so we can create the client here.
	return pb.NewServiceClient(conn), nil
}
```
