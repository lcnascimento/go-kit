package trace

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/lcnascimento/go-kit/env"
)

const defaultTraceRatio = 0.05

// RegisterFromEnv registers a global Tracer based on environment variables configuration.
// It returns the desired tracer and a shutdown function, along with an error if it happens.
// If errors are detected, we return a Noop tracer.
func RegisterFromEnv(ctx context.Context) (trace.Tracer, func(context.Context) error, error) {
	exporter, err := getSpanExporter(ctx)
	if err != nil {
		noopExporter := noop.NewTracerProvider().Tracer("")
		noopFlush := func(context.Context) error { return nil }

		return noopExporter, noopFlush, err
	}

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

	return tp.Tracer(serviceName), tp.Shutdown, nil
}

func getSpanExporter(ctx context.Context) (sdkTrace.SpanExporter, error) {
	switch strings.ToUpper(env.GetString("OTEL_SPAN_EXPORTER")) {
	case "OTLP":
		return getOTLPExporter(ctx)
	default:
		return nil, ErrSpanExporterNotSupported
	}
}
