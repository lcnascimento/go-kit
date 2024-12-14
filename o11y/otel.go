package o11y

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	iLog "github.com/lcnascimento/go-kit/o11y/log"
	"github.com/lcnascimento/go-kit/runtime"
)

var (
	pkg = "github.com/lcnascimento/go-kit/o11y"

	ctx context.Context

	tracerProvider *trace.TracerProvider
	meterProvider  *metric.MeterProvider
	loggerProvider *log.LoggerProvider

	logger *iLog.Logger
)

func init() {
	var err error

	ctx = runtime.ContextWithOSSignalCancellation()

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	tracerProvider, err = newTraceProvider(ctx)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tracerProvider)

	meterProvider, err = newMeterProvider(ctx)
	if err != nil {
		panic(err)
	}
	otel.SetMeterProvider(meterProvider)

	loggerProvider, err = newLoggerProvider(ctx)
	if err != nil {
		panic(err)
	}
	global.SetLoggerProvider(loggerProvider)

	logger = iLog.NewLogger(pkg)
	slog.SetDefault(slog.New(iLog.NewHandler(pkg)))
}

// Context returns the o11y context with OS signal cancellation.
func Context() context.Context {
	return ctx
}

// Shutdown shuts down the OpenTelemetry providers.
func Shutdown() {
	ctx := context.Background()

	if e := recover(); e != nil {
		logger.Criticalw(ctx, "panic", iLog.Any("exception", e))
	}

	if tracerProvider != nil {
		_ = tracerProvider.Shutdown(ctx)
	}

	if meterProvider != nil {
		_ = meterProvider.Shutdown(ctx)
	}

	if loggerProvider != nil {
		_ = loggerProvider.Shutdown(ctx)
	}
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	opts := []trace.TracerProviderOption{}

	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		otlpExporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient())
		if err != nil {
			return nil, err
		}

		opts = append(opts, trace.WithBatcher(otlpExporter))
	}

	return trace.NewTracerProvider(opts...), nil
}

func newMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	opts := []metric.Option{}

	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		otlpExporter, err := otlpmetricgrpc.New(ctx)
		if err != nil {
			return nil, err
		}

		opts = append(opts, metric.WithReader(metric.NewPeriodicReader(otlpExporter)))
	}

	return metric.NewMeterProvider(opts...), nil
}

func newLoggerProvider(ctx context.Context) (*log.LoggerProvider, error) {
	opts := []log.LoggerProviderOption{}

	logExporter, err := stdoutlog.New(stdoutlog.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	opts = append(opts, log.WithProcessor(log.NewBatchProcessor(logExporter)))

	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		otlpExporter, err := otlploggrpc.New(ctx)
		if err != nil {
			return nil, err
		}

		opts = append(opts, log.WithProcessor(log.NewBatchProcessor(otlpExporter)))
	}

	return log.NewLoggerProvider(opts...), nil
}
