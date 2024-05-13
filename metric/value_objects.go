package metric

import "go.opentelemetry.io/otel/metric"

// Meter is a wrapper around OTEL Meter.
type Meter interface {
	metric.Meter
}

// Aliases for commonly used features of OTEL metric package.
var (
	WithDescription  = metric.WithDescription
	WithUnit         = metric.WithUnit
	WithAttributeSet = metric.WithAttributeSet
	WithAttributes   = metric.WithAttributes
)
