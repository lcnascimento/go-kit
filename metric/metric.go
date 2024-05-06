package metric

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/lcnascimento/go-kit/env"
)

// RegisterFromEnv registers a global Meter based on environment variables configuration.
// It returns the desired meter and a shutdown function, along with an error if it happens.
// If errors are detected, we return a Noop meter.
func RegisterFromEnv(ctx context.Context) (metric.Meter, func(context.Context) error, error) {
	reader, err := getMetricReader(ctx)
	if err != nil {
		noopExporter := noop.NewMeterProvider().Meter("")
		noopFlush := func(context.Context) error { return nil }

		return noopExporter, noopFlush, err
	}

	serviceName := env.GetString("SERVICE_NAME", "unknown")
	serviceVersion := env.GetString("SERVICE_VERSION", "v0.0.0")

	opts := []sdkMetric.Option{
		sdkMetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			attribute.String("library.language", "go"),
		)),
		sdkMetric.WithReader(reader),
	}

	mp := sdkMetric.NewMeterProvider(opts...)

	otel.SetMeterProvider(mp)

	return mp.Meter(serviceName), mp.Shutdown, nil
}

func getMetricReader(ctx context.Context) (sdkMetric.Reader, error) {
	switch strings.ToUpper(env.GetString("OTEL_METRIC_READER")) {
	case "OTLP":
		return getOTLPReader(ctx)
	default:
		return nil, ErrMetricReaderNotSupported
	}
}
