package grpcserver

import (
	"context"
	"fmt"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/lcnascimento/go-kit/log"
)

const defaultPort = 3000

type server struct {
	mu *sync.Mutex

	app  string
	port int

	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor

	server   *grpc.Server
	listener net.Listener

	serviceRegistrations []ServiceRegistration
	services             []any
}

// NewServer creates a new gRPC server.
func NewServer(appName string, opts ...Option) Server {
	s := &server{
		mu:                   &sync.Mutex{},
		app:                  appName,
		port:                 defaultPort,
		unaryInterceptors:    []grpc.UnaryServerInterceptor{},
		streamInterceptors:   []grpc.StreamServerInterceptor{},
		serviceRegistrations: []ServiceRegistration{},
		services:             []any{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// RegisterService registers a service to the gRPC server.
func (s *server) RegisterService(registration ServiceRegistration, svc any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.serviceRegistrations = append(s.serviceRegistrations, registration)
	s.services = append(s.services, svc)
}

// Start starts the gRPC server.
func (s *server) Start(ctx context.Context) (err error) {
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		err = ErrCreateListener.WithCause(err)
		log.Error(ctx, err)
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
		service := s.services[i]

		register(s.server, service)
	}

	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()

	log.Info(ctx, "gRPC server started", log.Int("port", s.port))
	return s.server.Serve(s.listener)
}

// Stop stops the gRPC server.
func (s *server) Stop(ctx context.Context) error {
	if s.server == nil || s.listener == nil {
		return ErrServerNotStarted
	}

	log.Info(ctx, "stopping gRPC server")

	s.server.GracefulStop()
	return s.listener.Close()
}
