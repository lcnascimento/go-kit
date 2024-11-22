package main

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/log"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("example")
}

func main() {
	defer os.Exit(0)

	log.SetLevel(log.LevelDebug)

	ctx := context.Background()

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	defer log.Info(ctx, "Deferred")

	foo, _ := baggage.NewMember("foo", "foo")
	bar, _ := baggage.NewMember("bar", "bar")

	bag, _ := baggage.New(foo, bar)

	ctx = baggage.ContextWithBaggage(ctx, bag)

	attrs := []slog.Attr{
		log.String("attr1", "value1"),
		log.String("attr2", "value2"),
	}

	log.Debug(ctx, "Debug", attrs...)
	log.Info(ctx, "Info", attrs...)
	log.Warn(ctx, "Warn", attrs...)
	log.Errorw(ctx, "Error", attrs...)
	log.Criticalw(ctx, "Critical", attrs...)

	log.Error(ctx, errors.New("Custom Error"), attrs...)
	log.Critical(ctx, errors.New("Custom Critical Error"), attrs...)
}
