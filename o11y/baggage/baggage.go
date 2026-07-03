package baggage

import (
	"context"
	"net/url"

	oBaggage "go.opentelemetry.io/otel/baggage"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const MemberKeyCorrelationID = "correlation_id"

// type aliases for useful assets from the OTEL baggage package.
var FromContext = oBaggage.FromContext

// Member holds telemetry information in a key-value format.
type Member interface {
	Key() string
	Value() string
}

type member struct {
	key   string
	value string
}

func (m *member) Key() string {
	return m.key
}

func (m *member) Value() string {
	return m.value
}

// NewMember creates a new baggage member.
func NewMember(key, value string) Member {
	return &member{
		key:   key,
		value: url.QueryEscape(value),
	}
}

// ContextWithMembers returns a new context with the given members added to its baggage.
// If the baggage already has a member key, it will be replaced.
func ContextWithMembers(ctx context.Context, members ...Member) context.Context {
	bag := oBaggage.FromContext(ctx)
	span := trace.SpanFromContext(ctx)

	for _, m := range members {
		member, err := oBaggage.NewMember(m.Key(), m.Value())
		if err != nil {
			continue
		}

		bag, _ = bag.SetMember(member)

		if !span.IsRecording() {
			continue
		}

		span.SetAttributes(attribute.String(m.Key(), m.Value()))
	}

	return oBaggage.ContextWithBaggage(ctx, bag)
}

// ContextWithCorrelationID adds a correlation ID to the baggage.
// A CorrelationID differs from a trace ID. The TraceID is a unique identifier for a request.
// The CorrelationID is a unique identifier for a process, that can use more than one request to complete.
//
// Useful for external system calls.
func ContextWithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return ContextWithMembers(ctx, NewMember(MemberKeyCorrelationID, correlationID))
}
