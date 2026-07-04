//nolint:mnd // OK
package main

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y"
	"github.com/lcnascimento/go-kit/o11y/baggage"
	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	// This is how the OpenTelemetry community uses telemetry components.
	// Instantiate a tracer, meter and logger for each instrumented package, given them its exclusive name.
	// They do not create a single tracer, meter or logger for the entire application, passing it around via dependency injection.
	pkg    = "github.com/lcnascimento/go-kit/o11y/example"
	tracer = otel.Tracer(pkg)
	meter  = otel.Meter(pkg)
	logger = log.MustNewLogger(pkg)

	requestsCounter metric.Int64Counter
)

func main() {
	o11y.MustStart()
	defer o11y.Shutdown()

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch
		slog.Default().Debug("context canceled by external signal")
		cancel()
	}()

	requestsCounter, _ = meter.Int64Counter("o11y.example.requests")

	// To configure profiling, you can set the following environment variables:
	// OTEL_PROFILES_EXPORTER=pyroscope
	// OTEL_PROFILES=cpu,heap,goroutine,mutex
	// OTEL_PROFILE_EXPORT_INTERVAL=10s

	for {
		select {
		case <-time.After(time.Second):
			fakeRequest()
		case <-ctx.Done():
			return
		}
	}
}

//nolint:gosec // OK
func fakeRequest() {
	ctx, span := tracer.Start(context.Background(), "fake_request")
	defer span.End()

	ctx = baggage.ContextWithCorrelationID(ctx, "correlation-id")

	time.Sleep(time.Second)

	var code string
	switch rand.Intn(10) {
	case 0:
		code = "cached"
		logger.Debug(ctx, "nothing to do", slog.String("code", code))
	case 1:
		code = "warning"
		logger.Warn(ctx, "something might be wrong", slog.String("code", code))
	case 2:
		code = "error"
		logger.Error(ctx, fakeError("something is wrong", errors.CodeType(code)))
	case 3:
		code = "critical"
		logger.Critical(ctx, fakeError("something went really wrong", errors.CodeType(code)))
	default:
		code = "ok"
		logger.Info(ctx, "success", slog.String("code", code))
	}

	requestsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("code", code)))
}

func fakeError(msg string, code errors.CodeType) error {
	return errors.New("%s", msg).
		WithAttribute("key1", "value1").
		WithAttribute("key2", "value2").
		WithCode(code).
		WithKind(errors.KindResourceExhausted).
		WithCause(errors.New("nested error"))
}
