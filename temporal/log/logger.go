package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/workflow"

	otel "go.opentelemetry.io/otel/baggage"

	"github.com/lcnascimento/go-kit/o11y/log"

	"github.com/lcnascimento/go-kit/temporal/baggage"
	"github.com/lcnascimento/go-kit/temporal/internal/interceptors"
)

type (
	Attr  log.Attr
	Value log.Value
)

var (
	String     = log.String
	Int        = log.Int
	Float      = log.Float
	Bool       = log.Bool
	Any        = log.Any
	Group      = log.Group
	GroupValue = log.GroupValue
	Object     = log.Any
)

type Logger struct {
	logger log.Logger
}

func MustNewLogger(pkg string) *Logger {
	return &Logger{
		logger: log.MustNewLogger(pkg),
	}
}

func (l *Logger) Debug(ctx workflow.Context, msg string, attrs ...slog.Attr) {
	l.logger.Debug(build(ctx), msg, attrs...)
}

func (l *Logger) Info(ctx workflow.Context, msg string, attrs ...slog.Attr) {
	l.logger.Info(build(ctx), msg, attrs...)
}

func (l *Logger) Warn(ctx workflow.Context, msg string, attrs ...slog.Attr) {
	l.logger.Warn(build(ctx), msg, attrs...)
}

func (l *Logger) Error(ctx workflow.Context, err error, attrs ...slog.Attr) {
	l.logger.Error(build(ctx), err, attrs...)
}

func (l *Logger) ErrorMessage(ctx workflow.Context, msg string, attrs ...slog.Attr) {
	l.logger.ErrorMessage(build(ctx), msg, attrs...)
}

func (l *Logger) Critical(ctx workflow.Context, err error, attrs ...slog.Attr) {
	l.logger.Critical(build(ctx), err, attrs...)
}

func (l *Logger) CriticalMessage(ctx workflow.Context, msg string, attrs ...slog.Attr) {
	l.logger.CriticalMessage(build(ctx), msg, attrs...)
}

func (l *Logger) ErrorBySeverity(ctx workflow.Context, err error, attrs ...slog.Attr) {
	l.logger.ErrorBySeverity(build(ctx), err, attrs...)
}

func build(ctx workflow.Context) context.Context {
	out := context.Background()
	bag := baggage.FromContext(ctx)
	out = otel.ContextWithBaggage(out, bag)

	if span := interceptors.SpanFromWorkflowContext(ctx); span != nil {
		out = trace.ContextWithSpan(out, span)
	}

	return out
}
