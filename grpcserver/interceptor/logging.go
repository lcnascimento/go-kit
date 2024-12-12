package interceptor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/lcnascimento/go-kit/o11y/log"
)

const pkg = "github.com/lcnascimento/go-kit/grpcserver/interceptor"

var logger *log.Logger

func init() {
	logger = log.NewLogger(pkg)
}

const (
	// GRPCResponseLatencyKey is the amount of time needed to produce a response to a request.
	GRPCResponseLatencyKey = "rpc.grpc.response_latency"
)

// LoggingUnaryServerInterceptor returns a new unary interceptor suitable for request logging.
func LoggingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		now := time.Now()
		latency := now.Sub(start).String()

		attrs := []slog.Attr{
			slog.String(string(semconv.RPCMethodKey), info.FullMethod),
			slog.String(string(semconv.RPCGRPCStatusCodeKey), status.Code(err).String()),
			slog.String(GRPCResponseLatencyKey, latency),
		}

		peer, ok := peer.FromContext(ctx)
		if ok {
			attrs = append(attrs, slog.String(string(semconv.NetPeerIPKey), peer.Addr.String()))
		}

		msg := fmt.Sprintf("[RPC] %s", info.FullMethod)
		logger.Info(ctx, msg, attrs...)

		return resp, err
	}
}

// LoggingStreamServerInterceptor returns a new stream interceptor suitable for request logging.
func LoggingStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {

		ctx := ss.Context()
		start := time.Now()

		err = handler(srv, ss)

		now := time.Now()
		latency := now.Sub(start).String()

		attrs := []slog.Attr{
			slog.String(string(semconv.RPCMethodKey), info.FullMethod),
			slog.String(string(semconv.RPCGRPCStatusCodeKey), status.Code(err).String()),
			slog.String(GRPCResponseLatencyKey, latency),
		}

		peer, ok := peer.FromContext(ctx)
		if ok {
			attrs = append(attrs, slog.String(string(semconv.NetPeerIPKey), peer.Addr.String()))
		}

		msg := fmt.Sprintf("[RPC] %s", info.FullMethod)
		logger.Info(ctx, msg, attrs...)

		return nil
	}
}
