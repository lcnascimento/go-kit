package grpcserver

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Option func(*config)

func WithServerOpts(opts ...grpc.ServerOption) Option {
	return func(s *config) {
		s.serverOpts = append(s.serverOpts, opts...)
	}
}

func WithOtelOpts(opts ...otelgrpc.Option) Option {
	return func(s *config) {
		s.otelOpts = append(s.otelOpts, opts...)
	}
}
