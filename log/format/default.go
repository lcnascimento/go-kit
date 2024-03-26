package format

import (
	"context"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/errors"
)

// DefaultFormatter is a log formatter with very basic style.
type DefaultFormatter struct{}

// NewDefault creates a new default LogFormatter.
func NewDefault() *DefaultFormatter {
	return &DefaultFormatter{}
}

// Format formats the log payload that will be rendered.
func (b DefaultFormatter) Format(ctx context.Context, in *LogInput) any {
	payload := map[string]any{
		"level":     in.Level,
		"timestamp": in.Timestamp.Format(time.RFC3339),
		"message":   in.Message,
	}

	if in.Payload != nil {
		payload["payload"] = in.Payload
	}

	attrs := extractContextKeysFromContext(ctx, in.Attributes)
	if in.Err != nil {
		attrs[errors.ContextKeyRootError] = errors.RootError(in.Err)
		attrs[errors.ContextKeyErrorKind] = string(errors.Kind(in.Err))
		attrs[errors.ContextKeyErrorCode] = string(errors.Code(in.Err))
		attrs[errors.ContextKeyErrorRetryable] = strconv.FormatBool(errors.Retryable(in.Err))
	}

	if len(attrs) > 0 {
		payload["attributes"] = attrs
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().TraceID().IsValid() {
		return payload
	}

	span.AddEvent("log", trace.WithAttributes(buildOtelAttributes(attrs, "log")...))

	if in.Err != nil {
		span.RecordError(in.Err, trace.WithAttributes(buildOtelAttributes(attrs, "exception")...))
		span.SetStatus(codes.Error, in.Err.Error())
	}

	payload["trace_id"] = span.SpanContext().TraceID().String()
	payload["span_id"] = span.SpanContext().SpanID().String()

	return payload
}
