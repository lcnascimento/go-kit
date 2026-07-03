//nolint:revive // OK
package trace

import (
	"context"

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

	var exporter trace.SpanExporter
	switch cfg.TracesExporter {
	case config.TracesExporterOTLP:
		exporter, err = exporters.OtlpGRPC(ctx)
	case config.TracesExporterStdout:
		exporter, err = exporters.Stdout()
	default:
		return errors.New("invalid trace exporter: %s", cfg.TracesExporter)
	}

	if err != nil {
		return err
	}

	var sampler trace.Sampler
	switch cfg.TraceSampler {
	case config.TraceSamplerAlwaysOff:
		sampler = trace.NeverSample()
	case config.TraceSamplerAlwaysOn:
		sampler = trace.AlwaysSample()
	case config.TraceSamplerParentBasedAlwaysOn:
		sampler = trace.ParentBased(trace.AlwaysSample())
	case config.TraceSamplerParentBasedAlwaysOff:
		sampler = trace.ParentBased(trace.NeverSample())
	case config.TraceSamplerParentBasedTraceIDRatio:
		sampler = trace.ParentBased(trace.TraceIDRatioBased(cfg.TraceSamplerArg))
	}

	spanLimits := trace.SpanLimits{
		AttributeCountLimit:         int(cfg.SpanAttributeCountLimit),      //nolint:gosec // OK with overflow
		AttributePerEventCountLimit: int(cfg.SpanEventAttributeCountLimit), //nolint:gosec // OK with overflow
		AttributePerLinkCountLimit:  int(cfg.SpanLinkAttributeCountLimit),  //nolint:gosec // OK with overflow
		EventCountLimit:             int(cfg.SpanEventCountLimit),          //nolint:gosec // OK with overflow
		LinkCountLimit:              int(cfg.SpanLinkCountLimit),           //nolint:gosec // OK with overflow
	}

	attrValueLengthLimit := int(cfg.AttributeValueLengthLimit) //nolint:gosec // OK with overflow
	if cfg.SpanAttributeValueLengthLimit > 0 {
		attrValueLengthLimit = int(cfg.SpanAttributeValueLengthLimit) //nolint:gosec // OK with overflow
	}
	if attrValueLengthLimit > 0 {
		spanLimits.AttributeValueLengthLimit = attrValueLengthLimit
	}

	baggageProcessor := processors.BaggageSpanProcessor()

	processor := trace.NewBatchSpanProcessor(
		exporter,
		trace.WithBatchTimeout(cfg.BspExportTimeout),
		trace.WithMaxQueueSize(int(cfg.BspMaxQueueSize)),             //nolint:gosec // OK with overflow
		trace.WithMaxExportBatchSize(int(cfg.BspMaxExportBatchSize)), //nolint:gosec // OK with overflow
	)

	tp = trace.NewTracerProvider(
		trace.WithSampler(sampler),
		trace.WithSpanProcessor(baggageProcessor),
		trace.WithSpanProcessor(processor),
		trace.WithRawSpanLimits(spanLimits),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, cfg.ResourceAttributes.ToList()...)),
	)

	otel.SetTracerProvider(tp)
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
