package grpcserver

import "context"

// Server is a gRPC server.
type Server interface {
	// RegisterService registers a service to the gRPC server.
	RegisterService(registration ServiceRegistration)

	// Start starts the gRPC server.
	Start(context.Context) error

	// Stop stops the gRPC server.
	Stop(context.Context) error
}
