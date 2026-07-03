package grpcserver

import (
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/lcnascimento/go-kit/grpc/grpcserver/interceptor"
)

const (
	kb                      = 1024
	mb                      = kb * kb
	defaultFileMaxSize      = 20 * mb            // 20mb
	defaultGrpcMaxSize      = defaultFileMaxSize // should be the same size of file (upload size)
	defaultMaxGrpcMsgSize   = defaultGrpcMaxSize + mb
	defaultMinTime          = time.Second * 10
	defaultKeepAliveTimeout = time.Second * 10
	defaultPingInterval     = time.Second * 30
)

type config struct {
	serverOpts []grpc.ServerOption
	otelOpts   []otelgrpc.Option
}

// NewServer creates a new gRPC server with the given options.
func NewServer(opts ...Option) *grpc.Server {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	kaPolicy := keepalive.EnforcementPolicy{
		MinTime:             defaultMinTime, // min a client should wait to send a keepalive
		PermitWithoutStream: true,           // allow without ongoing streams
	}

	kaParams := keepalive.ServerParameters{
		Time:    defaultPingInterval,     // how often to ping
		Timeout: defaultKeepAliveTimeout, // how long will we wait for a reply
	}

	//nolint:prealloc // OK
	svrOpts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(defaultMaxGrpcMsgSize),
		grpc.KeepaliveEnforcementPolicy(kaPolicy),
		grpc.KeepaliveParams(kaParams),
		grpc.StatsHandler(otelgrpc.NewServerHandler(cfg.otelOpts...)),
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryLogging(),
			interceptor.UnaryErrorHandler(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.StreamLogging(),
			interceptor.StreamErrorHandler(),
		),
	}

	svrOpts = append(svrOpts, cfg.serverOpts...)
	return grpc.NewServer(svrOpts...)
}
