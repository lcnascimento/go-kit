package middleware

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// OpenTelemetry is a middleware that handles opentelemetry features.
// It assures context propagation between messages.
func OpenTelemetry() message.HandlerMiddleware {
	propagator := otel.GetTextMapPropagator()

	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			ctx := propagator.Extract(msg.Context(), propagation.MapCarrier(msg.Metadata))

			msg.SetContext(ctx)

			return h(msg)
		}
	}
}
