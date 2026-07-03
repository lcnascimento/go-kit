package interceptor

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/grpcserver/interceptor"
	logger = log.MustNewLogger(pkg)
)

// UnaryLogging returns a new unary interceptor suitable for request logging.
func UnaryLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		now := time.Now()
		latency := now.Sub(start).String()

		attrs := []log.Attr{
			log.String("code", status.Code(err).String()),
			log.String("latency", latency),
		}

		peer, ok := peer.FromContext(ctx)
		if ok {
			attrs = append(attrs, log.String("peer_ip", peer.Addr.String()))
		}

		msg := fmt.Sprintf("RPC %s", info.FullMethod)
		logger.Debug(ctx, msg, attrs...)

		return resp, err
	}
}

// StreamLogging returns a new stream interceptor suitable for request logging.
func StreamLogging() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		start := time.Now()

		err := handler(srv, ss)

		now := time.Now()
		latency := now.Sub(start).String()

		attrs := []log.Attr{
			log.String("code", status.Code(err).String()),
			log.String("latency", latency),
		}

		peer, ok := peer.FromContext(ctx)
		if ok {
			attrs = append(attrs, log.String("peer_ip", peer.Addr.String()))
		}

		msg := fmt.Sprintf("RPC %s", info.FullMethod)
		logger.Debug(ctx, msg, attrs...)

		return err
	}
}
