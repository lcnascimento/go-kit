package o11y

import (
	"context"
	"log/slog"

	oLog "go.opentelemetry.io/otel/log"
	oMetric "go.opentelemetry.io/otel/metric"
	oTrace "go.opentelemetry.io/otel/trace"

	pLog "github.com/lcnascimento/go-kit/o11y/log"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
	"github.com/lcnascimento/go-kit/o11y/internal/log"
	"github.com/lcnascimento/go-kit/o11y/internal/metric"
	"github.com/lcnascimento/go-kit/o11y/internal/profile"
	"github.com/lcnascimento/go-kit/o11y/internal/propagator"
	"github.com/lcnascimento/go-kit/o11y/internal/trace"
)

var (
	tp oTrace.TracerProvider
	mp oMetric.MeterProvider
	lp oLog.LoggerProvider
)

// Option is a function that configures the otel components.
type Option = pLog.HandlerOption

// type alias for handler options.
var (
	WithLogHandler         = pLog.WithLogHandler
	WithLoggerAttrResolver = pLog.WithLoggerAttrResolver
)

// Start starts the otel components.
func Start(opts ...Option) error {
	if tp != nil || mp != nil || lp != nil {
		return errors.New("otel already started")
	}

	ctx := context.Background()

	cfg := pLog.HandlerConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	if err := log.Start(ctx, cfg.Core(), cfg.AttrResolver()); err != nil {
		slog.Default().Error("could not start log component", "error", err)
		return err
	}
	if err := trace.Start(ctx); err != nil {
		slog.Default().Error("could not start trace component", "error", err)
		return err
	}
	if err := metric.Start(ctx); err != nil {
		slog.Default().Error("could not start metric component", "error", err)
		return err
	}
	if err := profile.Start(ctx); err != nil {
		slog.Default().Error("could not start profile component", "error", err)
		return err
	}

	propagator.Setup()

	slog.Default().Debug("o11y started successfully", "config", global.Config())
	return nil
}

// MustStart starts the otel components and panics if an error occurs.
func MustStart() {
	if err := Start(); err != nil {
		panic(err)
	}
}

// Shutdown shuts down the otel components.
func Shutdown() {
	ctx := context.Background()

	slog.Default().Debug("shutting down telemetry components")

	errs := []error{}

	if err := log.Shutdown(ctx); err != nil {
		errs = append(errs, err)
	}
	if err := metric.Shutdown(ctx); err != nil {
		errs = append(errs, err)
	}
	if err := trace.Shutdown(ctx); err != nil {
		errs = append(errs, err)
	}
	if err := profile.Shutdown(ctx); err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return
	}

	err := errors.New("could not shutdown telemetry components")
	for _, cause := range errs {
		err = err.WithCause(cause)
	}

	slog.Default().Error(err.Error(), "reasons", errors.Reasons(err))
}
