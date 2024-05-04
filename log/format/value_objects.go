package format

import (
	"context"
	"time"

	"github.com/lcnascimento/go-kit/propagation"
)

// AttributeSet defines custom attributes to be logged.
type AttributeSet map[string]any

// LogInput is the input given to a LogFormatter that is used to produce log payload.
type LogInput struct {
	Level       string
	Message     string
	Payload     any
	ContextKeys propagation.ContextKeySet
	Attributes  AttributeSet
	Timestamp   time.Time
}

func extractContextKeysFromContext(ctx context.Context, attrSet propagation.ContextKeySet) map[propagation.ContextKey]any {
	contextKeys := map[propagation.ContextKey]any{}

	for attr := range attrSet {
		if value := ctx.Value(attr); value != nil {
			contextKeys[attr] = value
		}
	}

	return contextKeys
}

func isError(level string) bool {
	return level == "ERROR" || level == "CRITICAL"
}
