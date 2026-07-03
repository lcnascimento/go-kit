package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/processors/minsev"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"

	oGlobal "go.opentelemetry.io/otel/log/global"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/config"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
	"github.com/lcnascimento/go-kit/o11y/internal/log/exporters"
)

const pkg = "github.com/lcnascimento/go-kit/o11y"

var lp *log.LoggerProvider

// Start starts the log component.
//
//nolint:gosec // OK with overflow
func Start(ctx context.Context, core slog.Handler, resolver AttrResolver) (err error) {
	cfg := global.Config()

	if cfg.Disabled {
		return nil
	}

	handler, err := NewHandler(pkg, core, resolver)
	if err != nil {
		return err
	}

	// Zap does not support Log Provider yet.
	if cfg.LogHandler == config.LogHandlerZap {
		cfg.LogsExporter = config.LogsExporterStdout
	}

	var exporter log.Exporter
	switch cfg.LogsExporter {
	case config.LogsExporterOTLP:
		exporter, err = exporters.OtlpGRPC(ctx)
	case config.LogsExporterStdout:
		exporter, err = exporters.Stdout()
	default:
		return errors.New("invalid log exporter: %s", cfg.LogsExporter)
	}

	if err != nil {
		return err
	}

	var processor log.Processor = log.NewBatchProcessor(
		exporter,
		log.WithExportTimeout(cfg.BlrpExportTimeout),
		log.WithExportMaxBatchSize(int(cfg.BlrpMaxExportBatchSize)),
		log.WithMaxQueueSize(int(cfg.BlrpMaxQueueSize)),
	)

	filterProcessor := minsev.NewLogProcessor(processor, cfg.LogLevel.ToSeverity())

	attrs := cfg.ResourceAttributes.ToList()
	attrs = append(attrs,
		semconv.ServiceNameKey.String(cfg.ServiceName),
		attribute.String("language", "go"),
	)

	opts := []log.LoggerProviderOption{
		log.WithProcessor(filterProcessor),
		log.WithResource(resource.NewWithAttributes(semconv.SchemaURL, attrs...)),
		log.WithAttributeCountLimit(int(cfg.LogRecordAttributeCountLimit)),
	}

	attrValueLengthLimit := int(cfg.AttributeValueLengthLimit)
	if cfg.LogRecordAttributeValueLengthLimit > 0 {
		attrValueLengthLimit = int(cfg.LogRecordAttributeValueLengthLimit)
	}
	if attrValueLengthLimit > 0 {
		opts = append(opts, log.WithAttributeValueLengthLimit(attrValueLengthLimit))
	}

	lp = log.NewLoggerProvider(opts...)

	slog.SetDefault(slog.New(handler))
	slog.SetLogLoggerLevel(cfg.LogLevel.ToSlogLevel())
	oGlobal.SetLoggerProvider(lp)

	return nil
}

// Shutdown shuts down the log component.
func Shutdown(ctx context.Context) error {
	if lp == nil {
		return nil
	}

	if err := lp.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown log provider")
	}

	return nil
}
