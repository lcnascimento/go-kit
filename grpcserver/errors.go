package grpcserver

import "github.com/lcnascimento/go-kit/errors"

var (
	// ErrCreateListener is returned when the listener creation fails.
	ErrCreateListener = errors.New("failed to create listener")

	// ErrServerNotStarted is returned when the server is not started yet.
	ErrServerNotStarted = errors.New("server not started")
)
