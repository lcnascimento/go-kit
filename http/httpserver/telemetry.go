package httpserver

import (
	"context"

	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/http/httpserver"
	logger = log.MustNewLogger(pkg)
)

func (s *Server) onStart(port int) {
	logger.Info(context.Background(), "starting HTTP server", log.Int("port", port))
}

func (s *Server) onShutdown(ctx context.Context) {
	logger.Info(ctx, "shutting down HTTP server")
}

func (s *Server) onError(ctx context.Context, err error) error {
	logger.ErrorBySeverity(ctx, err)

	return err
}
