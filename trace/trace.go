package trace

import (
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/lcnascimento/go-kit/env"
)

const defaultTraceRatio = 0.1

func init() {
	exporter, err := getExporter()
	if err != nil {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return
	}

	setupTracerProvider(exporter)
}

func getExporter() (sdkTrace.SpanExporter, error) {
	switch strings.ToUpper(env.GetString("OTEL_SPAN_EXPORTER")) {
	case "OTLP":
		return getOTLPExporter()
	default:
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	}
}

func setupTracerProvider(exporter sdkTrace.SpanExporter) {
	serviceName := env.GetString("SERVICE_NAME", "unknown")
	serviceVersion := env.GetString("SERVICE_VERSION", "v0.0.0")

	sampler := sdkTrace.TraceIDRatioBased(env.GetFloat("OTEL_TRACE_RATIO", defaultTraceRatio))

	options := []sdkTrace.TracerProviderOption{
		sdkTrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			attribute.String("library.language", "go"),
		)),
		sdkTrace.WithSampler(sdkTrace.ParentBased(sampler)),
		sdkTrace.WithBatcher(exporter),
	}

	tp := sdkTrace.NewTracerProvider(options...)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}
