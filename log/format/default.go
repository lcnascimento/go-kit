package format

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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

	contextKeys := extractContextKeysFromContext(ctx, in.ContextKeys)
	if len(contextKeys) > 0 {
		payload["context"] = contextKeys
	}

	if len(in.Attributes) > 0 {
		payload["attributes"] = in.Attributes
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().TraceID().IsValid() {
		return payload
	}

	if isError(in.Level) {
		span.SetStatus(codes.Error, in.Message)
	}

	payload["trace_id"] = span.SpanContext().TraceID().String()
	payload["span_id"] = span.SpanContext().SpanID().String()

	return payload
}
