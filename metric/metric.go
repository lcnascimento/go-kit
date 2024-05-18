package metric

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/lcnascimento/go-kit/env"
	"github.com/lcnascimento/go-kit/log"
)

var provider *metric.MeterProvider

func init() {
	exporter, err := getMetricExporter()
	if err != nil {
		log.Error(context.Background(), err)
		otel.SetMeterProvider(noop.NewMeterProvider())
		return
	}

	provider = setupMeterProvider(exporter)
}

// Shutdown flushes all metric data held by an exporter and releases any held computational resources.
func Shutdown(ctx context.Context) error {
	if provider == nil {
		return nil
	}

	log.Debug(ctx, "Shutting down Meter Provider...")
	return provider.Shutdown(ctx)
}

func getMetricExporter() (metric.Exporter, error) {
	ctx := context.Background()

	switch strings.ToUpper(env.GetString("OTEL_METRIC_EXPORTER")) {
	case "OTLP":
		log.Debug(ctx, "Installing OTLP Metric Exporter...")
		return getOTLPExporter()
	default:
		log.Debug(ctx, "Installing STDOUT Metric Exporter...")
		return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	}
}

func setupMeterProvider(exporter metric.Exporter) *metric.MeterProvider {
	serviceName := env.GetString("SERVICE_NAME", "unknown")
	serviceVersion := env.GetString("SERVICE_VERSION", "v0.0.0")

	reader := metric.NewPeriodicReader(exporter)

	opts := []metric.Option{
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			attribute.String("library.language", "go"),
		)),
		metric.WithReader(reader),
	}

	log.Debug(context.Background(), "Setting up Meter Provider...")
	mp := metric.NewMeterProvider(opts...)
	otel.SetMeterProvider(mp)

	return mp
}
