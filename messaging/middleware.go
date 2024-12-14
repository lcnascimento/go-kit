package messaging

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
)

var tracer = otel.Tracer("github.com/lcnascimento/go-kit/messaging")

// WithBaggage is a middleware that adds baggage to the message metadata.
func WithBaggage() message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			bag := baggage.FromContext(msg.Context())
			if bag.Len() > 0 {
				for _, member := range bag.Members() {
					msg.Metadata.Set(baggageFieldPrefix+member.Key(), member.Value())
				}
			}

			return h(msg)
		}
	}
}

// WithTracePropagation is a middleware that propagates the trace context to the message context.
func WithTracePropagation() message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			ctx := msg.Context()

			ctx, span := tracer.Start(ctx, msg.Metadata.Get("name"))
			defer span.End()

			msg.SetContext(ctx)

			return h(msg)
		}
	}
}
