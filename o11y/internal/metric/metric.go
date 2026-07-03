package metric

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/resource"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/config"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
	"github.com/lcnascimento/go-kit/o11y/internal/metric/exporters"
)

var mp *metric.MeterProvider

// Start starts metric provider.
func Start(ctx context.Context) (err error) {
	cfg := global.Config()

	if cfg.Disabled || cfg.MetricsExporter == config.MetricsExporterNone {
		return nil
	}

	var exporter metric.Exporter
	switch cfg.MetricsExporter {
	case config.MetricsExporterOTLP:
		exporter, err = exporters.OtlpGRPC(ctx)
	case config.MetricsExporterStdout:
		exporter, err = exporters.Stdout()
	default:
		return errors.New("invalid metrics exporter: %s", cfg.MetricsExporter)
	}

	if err != nil {
		return err
	}

	var filter exemplar.Filter
	switch cfg.MetricsExemplarFilter {
	case config.MetricsExemplarFilterAlwaysOn:
		filter = exemplar.AlwaysOnFilter
	case config.MetricsExemplarFilterAlwaysOff:
		filter = exemplar.AlwaysOffFilter
	case config.MetricsExemplarFilterTraceBased:
		filter = exemplar.TraceBasedFilter
	default:
		return errors.New("invalid metrics exemplar filter: %s", cfg.MetricsExemplarFilter)
	}

	reader := metric.NewPeriodicReader(
		exporter,
		metric.WithInterval(cfg.MetricExportInterval),
		metric.WithTimeout(cfg.MetricExportTimeout),
	)

	mp = metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithExemplarFilter(filter),
		metric.WithResource(resource.NewWithAttributes(semconv.SchemaURL, cfg.ResourceAttributes.ToList()...)),
	)

	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		return errors.Wrap(err, "failed to start runtime metrics")
	}

	otel.SetMeterProvider(mp)
	return nil
}

// Shutdown shuts down the metric provider.
func Shutdown(ctx context.Context) error {
	if mp == nil {
		return nil
	}

	if err := mp.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown metrics provider")
	}

	return nil
}
