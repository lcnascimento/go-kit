package exporters

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

const (
	reconnectionPeriod           = 2 * time.Second
	defaultDialBackoffBaseDelay  = 1 * time.Second
	defaultDialBackoffMultiplier = 1.6
	defaultDialBackoffMaxDelay   = 15 * time.Second
	defaultDialMinConnectTimeout = 0
	defaultMaxMsgSize            = 1024 * 1024 * 5 // 5MB
)

// OtlpGRPC creates a new OTLP gRPC exporter.
func OtlpGRPC(ctx context.Context) (trace.SpanExporter, error) {
	cfg := global.Config()

	endpoint := cfg.OtlpTracesExporterEndpoint
	if endpoint == "" {
		endpoint = cfg.OtlpExporterEndpoint
	}

	connect := grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  defaultDialBackoffBaseDelay,
			Multiplier: defaultDialBackoffMultiplier,
			MaxDelay:   defaultDialBackoffMaxDelay,
		},
		MinConnectTimeout: defaultDialMinConnectTimeout,
	}

	dial := []grpc.DialOption{
		grpc.WithConnectParams(connect),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaultMaxMsgSize)),
	}

	clientOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithReconnectionPeriod(reconnectionPeriod),
		otlptracegrpc.WithTimeout(cfg.BspExportTimeout),
		otlptracegrpc.WithCompressor("gzip"),
		otlptracegrpc.WithDialOption(dial...),
	}

	if cfg.OtlpTracesInsecure {
		clientOpts = append(clientOpts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(clientOpts...))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create otlp trace grpc exporter")
	}

	return exporter, nil
}
