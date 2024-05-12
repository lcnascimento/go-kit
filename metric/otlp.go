package metric

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/lcnascimento/go-kit/env"
)

const (
	otlpDefaultEndpoint          = "http://localhost:4317"
	otlpDefaultReconnectPeriod   = 2
	otlpDefaultTimeout           = 30
	otlpDefaultBackoffBaseDelay  = 1
	otlpDefaultBackoffMaxDelay   = 15
	otlpDefaultBackoffMultiplier = 1.6
)

func getOTLPExporter() (metric.Exporter, error) {
	endpoint := env.GetString("OTEL_OTLP_ENDPOINT", otlpDefaultEndpoint)
	reconnectPeriod := env.GetInt("OTEL_OTLP_RECONNECT_PERIOD_IN_SECONDS", otlpDefaultReconnectPeriod)
	timeout := env.GetInt("OTEL_OTLP_TIMEOUT_IN_SECONDS", otlpDefaultTimeout)

	backoffBaseDelay := env.GetInt("OTEL_OTLP_BACKOFF_BASE_DELAY_IN_SECONDS", otlpDefaultBackoffBaseDelay)
	backoffMaxDelay := env.GetInt("OTEL_OTLP_BACKOFF_MAX_DELAY_IN_SECONDS", otlpDefaultBackoffMaxDelay)
	backoffMultiplier := env.GetFloat("OTEL_OTLP_BACKOFF_MULTIPLIER", otlpDefaultBackoffMultiplier)

	clientOpts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithReconnectionPeriod(time.Second * time.Duration(reconnectPeriod)),
		otlpmetricgrpc.WithTimeout(time.Second * time.Duration(timeout)),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithDialOption(
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  time.Second * time.Duration(backoffBaseDelay),
					Multiplier: backoffMultiplier,
					MaxDelay:   time.Second * time.Duration(backoffMaxDelay),
				},
				MinConnectTimeout: 0,
			}),
		),
	}

	return otlpmetricgrpc.New(context.Background(), clientOpts...)
}
