package trace

import (
	"context"

	otelpyroscope "github.com/grafana/otel-profiling-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/config"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
	"github.com/lcnascimento/go-kit/o11y/internal/trace/exporters"
	"github.com/lcnascimento/go-kit/o11y/internal/trace/processors"
)

var tp *trace.TracerProvider

// Start starts the trace provider.
func Start(ctx context.Context) (err error) {
	cfg := global.Config()

	if cfg.Disabled || cfg.TracesExporter == config.TracesExporterNone {
		return nil
	}

	exporter, err := newExporter(ctx, cfg.TracesExporter)
	if err != nil {
		return err
	}

	spanLimits := trace.SpanLimits{
		AttributeCountLimit:         int(cfg.SpanAttributeCountLimit),
		AttributePerEventCountLimit: int(cfg.SpanEventAttributeCountLimit),
		AttributePerLinkCountLimit:  int(cfg.SpanLinkAttributeCountLimit),
		EventCountLimit:             int(cfg.SpanEventCountLimit),
		LinkCountLimit:              int(cfg.SpanLinkCountLimit),
	}

	attrValueLengthLimit := int(cfg.AttributeValueLengthLimit)
	if cfg.SpanAttributeValueLengthLimit > 0 {
		attrValueLengthLimit = int(cfg.SpanAttributeValueLengthLimit)
	}
	if attrValueLengthLimit > 0 {
		spanLimits.AttributeValueLengthLimit = attrValueLengthLimit
	}

	baggageProcessor := processors.BaggageSpanProcessor()

	processor := trace.NewBatchSpanProcessor(
		exporter,
		trace.WithBatchTimeout(cfg.BspExportTimeout),
		trace.WithMaxQueueSize(int(cfg.BspMaxQueueSize)),
		trace.WithMaxExportBatchSize(int(cfg.BspMaxExportBatchSize)),
	)

	tp = trace.NewTracerProvider(
		trace.WithSampler(newSampler(cfg.TraceSampler, cfg.TraceSamplerArg)),
		trace.WithSpanProcessor(baggageProcessor),
		trace.WithSpanProcessor(processor),
		trace.WithRawSpanLimits(spanLimits),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, cfg.ResourceAttributes.ToList()...)),
	)

	if cfg.ProfilesExporter == config.ProfilesExporterNone {
		otel.SetTracerProvider(tp)
		return nil
	}

	otel.SetTracerProvider(otelpyroscope.NewTracerProvider(tp))
	return nil
}

// Shutdown shuts down the trace provider.
func Shutdown(ctx context.Context) error {
	if tp == nil {
		return nil
	}

	if err := tp.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown trace provider")
	}

	return nil
}

func newExporter(ctx context.Context, kind config.TracesExporter) (trace.SpanExporter, error) {
	switch kind {
	case config.TracesExporterOTLP:
		return exporters.OtlpGRPC(ctx)
	case config.TracesExporterStdout:
		return exporters.Stdout()
	default:
		return nil, errors.New("invalid trace exporter: %s", kind)
	}
}

func newSampler(kind config.TraceSampler, ratio float64) trace.Sampler {
	switch kind {
	case config.TraceSamplerAlwaysOff:
		return trace.NeverSample()
	case config.TraceSamplerAlwaysOn:
		return trace.AlwaysSample()
	case config.TraceSamplerParentBasedAlwaysOn:
		return trace.ParentBased(trace.AlwaysSample())
	case config.TraceSamplerParentBasedAlwaysOff:
		return trace.ParentBased(trace.NeverSample())
	case config.TraceSamplerParentBasedTraceIDRatio:
		return trace.ParentBased(trace.TraceIDRatioBased(ratio))
	default:
		// nil keeps the SDK default.
		return nil
	}
}
