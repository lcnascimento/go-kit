package metric

import (
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/lcnascimento/go-kit/env"
)

func init() {
	exporter, err := getMetricExporter()
	if err != nil {
		otel.SetMeterProvider(noop.NewMeterProvider())
		return
	}

	setupMeterProvider(exporter)
}

func getMetricExporter() (metric.Exporter, error) {
	switch strings.ToUpper(env.GetString("OTEL_METRIC_EXPORTER")) {
	case "OTLP":
		return getOTLPExporter()
	default:
		return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	}
}

func setupMeterProvider(exporter metric.Exporter) {
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

	mp := metric.NewMeterProvider(opts...)

	otel.SetMeterProvider(mp)
}
