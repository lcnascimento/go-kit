package propagator

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/lcnascimento/go-kit/o11y/internal/config"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

// Setup sets up the propagator.
func Setup() {
	cfg := global.Config()

	propagators := []propagation.TextMapPropagator{}
	if cfg.Propagators != nil {
		for prop := range cfg.Propagators {
			switch prop {
			case config.PropagatorTraceContext:
				propagators = append(propagators, propagation.TraceContext{})
			case config.PropagatorBaggage:
				propagators = append(propagators, propagation.Baggage{})
			}
		}
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagators...))
}
