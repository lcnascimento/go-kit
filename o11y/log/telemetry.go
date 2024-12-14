package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/errors"
)

const (
	pkg               = "github.com/lcnascimento/go-kit/o11y/log"
	errorsCounterName = "errors_count"
)

var errorsCounter metric.Int64Counter

func init() {
	var err error

	meter := otel.Meter(pkg)

	errorsCounter, err = meter.Int64Counter(errorsCounterName)
	if err != nil {
		panic(err)
	}
}

func onError(ctx context.Context, level slog.Level, err error) error {
	kvs := []attribute.KeyValue{
		attribute.String("code", string(errors.Code(err))),
		attribute.String("kind", string(errors.Kind(err))),
		attribute.String("severity", levelNames[level]),
		attribute.Bool("retryable", errors.IsRetryable(err)),
	}

	errorsCounter.Add(ctx, 1, metric.WithAttributes(kvs...))
	errorsCounter.Add(ctx, 1, metric.WithAttributes(kvs...))

	causes := make([]string, 0)
	for _, e := range errors.Unwrap(err) {
		causes = append(causes, e.Error())
	}
	kvs = append(kvs, attribute.StringSlice("causes", causes))

	span := trace.SpanFromContext(ctx)
	span.RecordError(err, trace.WithAttributes(kvs...), trace.WithStackTrace(true))
	span.SetStatus(codes.Error, err.Error())

	return err
}
