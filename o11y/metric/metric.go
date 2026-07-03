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

func MustFloat64ObservableGauge(meter metric.Meter, name, description string) metric.Float64ObservableGauge {
	gauge, err := meter.Float64ObservableGauge(name, metric.WithDescription(description), metric.WithUnit("1"))
	if err != nil {
		panic(err)
	}

	return gauge
}
