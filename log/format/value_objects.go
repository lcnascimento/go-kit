package format

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"github.com/lcnascimento/go-kit/propagation"
)

// LogInput is the input given to a LogFormatter that is used to produce log payload.
type LogInput struct {
	Level      string
	Message    string
	Err        error
	Payload    any
	Attributes propagation.ContextKeySet
	Timestamp  time.Time
}

func extractContextKeysFromContext(ctx context.Context, attrSet propagation.ContextKeySet) map[propagation.ContextKey]any {
	attributes := map[propagation.ContextKey]any{}

	for attr := range attrSet {
		if value := ctx.Value(attr); value != nil {
			attributes[attr] = value
		}
	}

	return attributes
}

func buildOtelAttributes(attrs map[propagation.ContextKey]any, prefix string) []attribute.KeyValue {
	eAttrs := []attribute.KeyValue{}

	for k, v := range attrs {
		key := fmt.Sprintf("%s.%s", prefix, k)
		value := fmt.Sprintf("%v", v)

		eAttrs = append(eAttrs, attribute.String(key, value))
	}

	return eAttrs
}
