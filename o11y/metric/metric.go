package metric

import "go.opentelemetry.io/otel/metric"

var (
	WithAttributes   = metric.WithAttributes
	WithAttributeSet = metric.WithAttributeSet
)

func MustIntCounter(meter metric.Meter, name, description string) metric.Int64Counter {
	counter, err := meter.Int64Counter(name, metric.WithDescription(description), metric.WithUnit("1"))
	if err != nil {
		panic(err)
	}

	return counter
}

func MustUpDownCounter(meter metric.Meter, name, description string) metric.Int64UpDownCounter {
	counter, err := meter.Int64UpDownCounter(name, metric.WithDescription(description), metric.WithUnit("1"))
	if err != nil {
		panic(err)
	}

	return counter
}

func MustFloat64Histogram(meter metric.Meter, name, description string) metric.Float64Histogram {
	histogram, err := meter.Float64Histogram(name, metric.WithDescription(description), metric.WithUnit("ms"))
	if err != nil {
		panic(err)
	}

	return histogram
}

var semconvDurationBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10}

// MustDurationHistogram creates a histogram in seconds with the OpenTelemetry semantic convention duration buckets.
func MustDurationHistogram(meter metric.Meter, name, description string) metric.Float64Histogram {
	histogram, err := meter.Float64Histogram(name,
		metric.WithDescription(description),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(semconvDurationBuckets...))
	if err != nil {
		panic(err)
	}

	return histogram
}

func MustFloat64ObservableGauge(meter metric.Meter, name, description string) metric.Float64ObservableGauge {
	gauge, err := meter.Float64ObservableGauge(name, metric.WithDescription(description), metric.WithUnit("1"))
	if err != nil {
		panic(err)
	}

	return gauge
}
