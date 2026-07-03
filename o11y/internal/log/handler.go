package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/config"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
	"github.com/lcnascimento/go-kit/o11y/internal/log/handlers"
)

// AttrResolver is a function that resolves attributes for a log record.
type AttrResolver func(context.Context, slog.Record) []slog.Attr

type handler struct {
	core     slog.Handler
	resolver AttrResolver
}

// NewHandler creates a new slog.Handler with the given options.
func NewHandler(name string, core slog.Handler, resolver AttrResolver) (slog.Handler, error) {
	cfg := global.Config()

	h := &handler{core: core, resolver: resolver}

	if h.core != nil {
		return h, nil
	}

	var err error

	switch cfg.LogHandler {
	case config.LogHandlerZap:
		h.core, err = handlers.Zap(name)
	case config.LogHandlerOtelSlog:
		h.core = otelslog.NewHandler(name)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to create slog handler")
	}

	return h, nil
}

// Enabled reports whether the handler handles records at the given level. The handler ignores records whose level is lower.
// It is called early, before any arguments are processed, to save effort if the log event should be discarded.
// If called from a Logger method, the first argument is the context passed to that method, or context.Background()
// if nil was passed or the method does not take a context. The context is passed so Enabled can use its values to make a decision.
func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.core.Enabled(ctx, level)
}

// Handle handles the Record. It resolves the attributes for the log record and then passes it to the core handler.
//
//nolint:gocritic // OK.
func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	attrs := []slog.Attr{}
	if bag := baggageAttr(ctx); bag.Key != "" {
		attrs = append(attrs, bag)
	}
	record.AddAttrs(attrs...)

	if h.resolver != nil {
		attrs = h.resolver(ctx, record)
		record.AddAttrs(attrs...)
	}

	return h.core.Handle(ctx, record)
}

// WithAttrs returns a new Handler whose attributes consist of both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{core: h.core.WithAttrs(attrs), resolver: h.resolver}
}

// WithGroup returns a new Handler with the given group appended to the receiver's existing groups.
// The keys of all subsequent attributes, whether added by With or in a Record, should be qualified by the sequence of group names.
func (h *handler) WithGroup(group string) slog.Handler {
	return &handler{core: h.core.WithGroup(group), resolver: h.resolver}
}
