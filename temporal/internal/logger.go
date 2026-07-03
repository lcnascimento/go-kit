package internal

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/log"
	"github.com/lcnascimento/go-kit/util"
)

var logger = log.MustNewLogger("github.com/lcnascimento/go-kit/temporal")

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Debug(msg string, args ...any) {
	ctx, attrs := buildLogData(args)
	logger.Debug(ctx, msg, attrs...)
}

func (l *Logger) Info(msg string, args ...any) {
	ctx, attrs := buildLogData(args)
	logger.Info(ctx, msg, attrs...)
}

func (l *Logger) Warn(msg string, args ...any) {
	ctx, attrs := buildLogData(args)
	logger.Warn(ctx, msg, attrs...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	_, attrs := buildLogData(args)
	logger.Warn(ctx, msg, attrs...)
}

func (l *Logger) Error(msg string, args ...any) {
	const actErrMsg = "Activity error."
	if msg == actErrMsg { // use ActivityError interceptor instead
		return
	}

	ctx, attrs := buildLogData(args)
	err := errors.New("%s", msg)

	logger.ErrorBySeverity(ctx, err, attrs...)
}

func buildLogData(args []any) (ctx context.Context, attrs []log.Attr) {
	ctx = context.Background()
	attrs = []log.Attr{}

	if len(args)%2 != 0 {
		return ctx, attrs
	}

	var (
		traceID trace.TraceID
		spanID  trace.SpanID
	)

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		value := args[i+1]

		if key == "TraceID" {
			traceID, _ = value.(trace.TraceID)

			continue
		}

		if key == "SpanID" {
			spanID, _ = value.(trace.SpanID)

			continue
		}

		if temporalKeys[key] {
			key = "temporal." + util.ToSnakeCase(key)
		}

		attrs = append(attrs, log.Any(key, value))
	}

	if traceID != [16]byte{} && spanID != [8]byte{} {
		ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
		}))
	}

	return ctx, attrs
}

var temporalKeys = map[string]bool{
	"Namespace":    true,
	"TaskQueue":    true,
	"WorkerID":     true,
	"WorkflowID":   true,
	"WorkflowType": true,
	"ActivityType": true,
	"RunID":        true,
	"Attempt":      true,
	"Error":        true,
}
