package exporters

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/sdk/log"
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

// OtlpGRPC creates a new otlp grpc exporter.
func OtlpGRPC(ctx context.Context) (log.Exporter, error) {
	cfg := global.Config()

	endpoint := cfg.OtlpLogsExporterEndpoint
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

	clientOpts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithReconnectionPeriod(reconnectionPeriod),
		otlploggrpc.WithTimeout(cfg.BlrpExportTimeout),
		otlploggrpc.WithCompressor("gzip"),
		otlploggrpc.WithDialOption(dial...),
	}

	if cfg.OtlpTracesInsecure {
		clientOpts = append(clientOpts, otlploggrpc.WithInsecure())
	}

	exporter, err := otlploggrpc.New(ctx, clientOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create otlp log grpc exporter")
	}

	return exporter, nil
}
