package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/lcnascimento/go-kit/env"

	"github.com/lcnascimento/go-kit/http/httpserver/middlewares"
)

const defaultReadHeaderTimeout = 5 * time.Second

type Server struct {
	server *http.Server
}

func NewServer(opts ...Option) *Server {
	svr := &Server{
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", env.Get("PORT", env.WithDefaultValue(3000))),
			ReadHeaderTimeout: defaultReadHeaderTimeout,
		},
	}

	for _, opt := range opts {
		opt(svr)
	}

	return svr
}

func (s *Server) Start(cb func(router *mux.Router) error) error {
	router := mux.NewRouter()

	router.StrictSlash(true)

	router.Use(middlewares.CorrelationID)
	router.Use(middlewares.Telemetry)
	router.Use(middlewares.Recover)

	if err := cb(router); err != nil {
		return err
	}

	s.server.Handler = router

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
