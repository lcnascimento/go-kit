package log

import (
	"context"
	"log/slog"

	"github.com/lcnascimento/go-kit/o11y/internal/log"
)

// HandlerConfig holds the configuration for the slog.Handler.
type HandlerConfig struct {
	core     slog.Handler
	resolver func(context.Context, slog.Record) []slog.Attr
}

// NewHandler creates a new slog.Handler with the given options.
func NewHandler(name string, opts ...HandlerOption) (slog.Handler, error) {
	h := &HandlerConfig{}
	for _, opt := range opts {
		opt(h)
	}

	return log.NewHandler(name, h.core, h.resolver)
}

// Core returns the core slog.Handler.
func (h *HandlerConfig) Core() slog.Handler {
	return h.core
}

// AttrResolver returns the attribute resolver.
func (h *HandlerConfig) AttrResolver() func(context.Context, slog.Record) []slog.Attr {
	return h.resolver
}

// HandlerOption is a function that configures a Handler.
type HandlerOption func(*HandlerConfig)

// WithLogHandler sets the core handler for the handler.
func WithLogHandler(core slog.Handler) HandlerOption {
	return func(h *HandlerConfig) {
		h.core = core
	}
}

// WithLoggerAttrResolver sets the logger attribute resolver in the config.
func WithLoggerAttrResolver(resolver func(context.Context, slog.Record) []slog.Attr) HandlerOption {
	return func(h *HandlerConfig) {
		h.resolver = resolver
	}
}
