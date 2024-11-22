package trace

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"

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

func getOTLPExporter() (trace.SpanExporter, error) {
	endpoint := env.GetString("OTEL_OTLP_ENDPOINT", otlpDefaultEndpoint)
	reconnectPeriod := env.GetInt("OTEL_OTLP_RECONNECT_PERIOD_IN_SECONDS", otlpDefaultReconnectPeriod)
	timeout := env.GetInt("OTEL_OTLP_TIMEOUT_IN_SECONDS", otlpDefaultTimeout)

	backoffBaseDelay := env.GetInt("OTEL_OTLP_BACKOFF_BASE_DELAY_IN_SECONDS", otlpDefaultBackoffBaseDelay)
	backoffMaxDelay := env.GetInt("OTEL_OTLP_BACKOFF_MAX_DELAY_IN_SECONDS", otlpDefaultBackoffMaxDelay)
	backoffMultiplier := env.GetFloat("OTEL_OTLP_BACKOFF_MULTIPLIER", otlpDefaultBackoffMultiplier)

	clientOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithReconnectionPeriod(time.Second * time.Duration(reconnectPeriod)),
		otlptracegrpc.WithTimeout(time.Second * time.Duration(timeout)),
		otlptracegrpc.WithCompressor("gzip"),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(
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

	ctx := context.Background()
	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(clientOpts...))
	if err != nil {
		return nil, err
	}

	return exporter, nil
}
