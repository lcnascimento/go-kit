package grpcserver

import (
	"context"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

const defaultPort = 3000

type server struct {
	app  string
	port int

	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor

	server   *grpc.Server
	listener net.Listener

	serviceRegistrations []ServiceRegistration
}

// NewServer creates a new gRPC server.
func NewServer(appName string, opts ...Option) Server {
	s := &server{
		app:                  appName,
		port:                 defaultPort,
		unaryInterceptors:    []grpc.UnaryServerInterceptor{},
		streamInterceptors:   []grpc.StreamServerInterceptor{},
		serviceRegistrations: []ServiceRegistration{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// RegisterService registers a service to the gRPC server.
func (s *server) RegisterService(registration ServiceRegistration) {
	s.serviceRegistrations = append(s.serviceRegistrations, registration)
}

// Start starts the gRPC server.
func (s *server) Start(ctx context.Context) (err error) {
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		onCreateListenerError(ctx, err)
		return err
	}

	options := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(s.unaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(s.streamInterceptors...)),
	}

	s.server = grpc.NewServer(options...)

	for i := range s.serviceRegistrations {
		register := s.serviceRegistrations[i]
		register(s.server)
	}

	done := make(chan error)
	defer close(done)

	go func() {
		<-ctx.Done()
		done <- nil
	}()

	go func() {
		onStart(ctx, s.port)
		done <- s.server.Serve(s.listener)
	}()

	return <-done
}

// Stop stops the gRPC server.
func (s *server) Stop(ctx context.Context) error {
	if s.server == nil || s.listener == nil {
		return ErrServerNotStarted
	}

	onStop(ctx)

	s.server.GracefulStop()
	return s.listener.Close()
}
