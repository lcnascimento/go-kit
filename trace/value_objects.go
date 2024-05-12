package trace

import "go.opentelemetry.io/otel/trace"

// Tracer is a wrapper around OTEL Tracer.
type Tracer interface {
	trace.Tracer
}
