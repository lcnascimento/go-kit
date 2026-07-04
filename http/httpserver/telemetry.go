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

func (s *Server) onShutdown() {
	logger.Info(context.Background(), "shutting down HTTP server")
}
