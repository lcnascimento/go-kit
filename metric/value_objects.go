package metric

import "go.opentelemetry.io/otel/metric"

// Meter is a wrapper around OTEL Meter.
type Meter interface {
	metric.Meter
}
