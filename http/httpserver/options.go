package httpserver

import (
	"fmt"
	"time"
)

type Option func(*Server)

func WithPort(port int) Option {
	return func(s *Server) {
		s.server.Addr = fmt.Sprintf(":%d", port)
	}
}

func WithReadHeaderTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadHeaderTimeout = timeout
	}
}
