package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"

	"github.com/lcnascimento/go-kit/o11y"
	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/o11y/example"
	tracer = otel.Tracer(pkg)
	logger = log.NewLogger(pkg)
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	defer o11y.Shutdown(ctx)

	foo, _ := baggage.NewMember("foo", "foo")
	bar, _ := baggage.NewMember("bar", "bar")

	bag, _ := baggage.New(foo, bar)

	ctx = baggage.ContextWithBaggage(ctx, bag)

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	attrs := []slog.Attr{
		slog.String("yay", "keke"),
	}

	logger.Debug(ctx, "DEBUG", attrs...)
	logger.Info(ctx, "INFO", attrs...)
	logger.Warn(ctx, "WARN", attrs...)
	logger.Error(ctx, errors.New("ERROR"), attrs...)
	logger.Critical(ctx, errors.New("CRITICAL"), attrs...)
}
