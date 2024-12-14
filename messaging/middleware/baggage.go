package middleware

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/baggage"
)

// WithBaggage is a middleware that adds baggage to the message metadata.
func WithBaggage(prefix string) message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			bag := baggage.FromContext(msg.Context())
			if bag.Len() > 0 {
				for _, member := range bag.Members() {
					msg.Metadata.Set(prefix+member.Key(), member.Value())
				}
			}

			return h(msg)
		}
	}
}
