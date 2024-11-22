package grpcserver

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"

	"github.com/lcnascimento/go-kit/grpcserver/interceptor"
)

// ServiceRegistration is a function that registers a service to the gRPC server.
type ServiceRegistration func(grpc.ServiceRegistrar, any)

// Option is a function that configures the gRPC server.
type Option func(*server)

// WithPort sets the port for the gRPC server.
func WithPort(port int) Option {
	return func(s *server) {
		s.port = port
	}
}

// WithUnaryInterceptor adds a unary interceptor to the gRPC server.
func WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) Option {
	return func(s *server) {
		s.unaryInterceptors = append(s.unaryInterceptors, interceptor)
	}
}

// WithStreamInterceptor adds a stream interceptor to the gRPC server.
func WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) Option {
	return func(s *server) {
		s.streamInterceptors = append(s.streamInterceptors, interceptor)
	}
}

// WithDefaultInterceptors adds the default interceptors to the gRPC server.
func WithDefaultInterceptors() Option {
	return func(s *server) {
		s.unaryInterceptors = append(s.unaryInterceptors, []grpc.UnaryServerInterceptor{
			interceptor.LoggingUnaryServerInterceptor(),
			interceptor.ErrorHandlingUnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		}...)
		s.streamInterceptors = append(s.streamInterceptors, []grpc.StreamServerInterceptor{
			interceptor.LoggingStreamServerInterceptor(),
			interceptor.ErrorHandlingStreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		}...)
	}
}
