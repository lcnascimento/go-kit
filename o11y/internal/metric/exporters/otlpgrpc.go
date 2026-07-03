package exporters

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

const (
	defaultReconnectionPeriod    = 2 * time.Second
	defaultDialBackoffBaseDelay  = 1 * time.Second
	defaultDialBackoffMultiplier = 1.6
	defaultDialBackoffMaxDelay   = 15 * time.Second
	defaultDialMinConnectTimeout = 0
	defaultMaxMsgSize            = 1024 * 1024 * 5 // 5MB
)

// OtlpGRPC creates a new OTLP gRPC metric exporter.
func OtlpGRPC(ctx context.Context) (metric.Exporter, error) {
	cfg := global.Config()

	endpoint := cfg.OtlpMetricsExporterEndpoint
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

	client := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithTimeout(cfg.MetricExportTimeout),
		otlpmetricgrpc.WithReconnectionPeriod(defaultReconnectionPeriod),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithDialOption(dial...),
	}

	if cfg.OtlpInsecure || cfg.OtlpMetricsInsecure {
		client = append(client, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, client...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create otlp metric grpc exporter")
	}

	return exporter, nil
}
