package grpcserver

import (
	"context"

	"github.com/lcnascimento/go-kit/o11y/log"
)

const pkg = "github.com/lcnascimento/go-kit/grpcserver"

var logger *log.Logger

func init() {
	logger = log.NewLogger(pkg)
}

func onCreateListenerError(ctx context.Context, err error) {
	err = ErrCreateListener.WithCause(err)
	logger.Error(ctx, err)
}

func onStart(ctx context.Context, port int) {
	logger.Info(ctx, "gRPC server started", log.Int("port", port))
}

func onStop(ctx context.Context) {
	logger.Info(ctx, "gRPC server stopped")
}
