package main

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"
	"github.com/lcnascimento/go-kit/kafka"
	"github.com/lcnascimento/go-kit/kafka/example"
	"github.com/lcnascimento/go-kit/o11y"
	"github.com/lcnascimento/go-kit/o11y/baggage"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/kafka/example/producer"
	tracer = otel.Tracer(pkg)
)

func main() {
	ctx := context.Background()
	o11y.MustStart(ctx)
	defer o11y.Shutdown(ctx)

	producer := kafka.NewProducer()
	defer producer.Stop(ctx)

	ctx, span := tracer.Start(ctx, "ExampleProducer")
	defer span.End()

	ctx = baggage.ContextWithCorrelationID(ctx, uuid.Must(uuid.NewV7()).String())
	ctx = baggage.ContextWithMembers(ctx, baggage.NewMember("foo", "bar"))

	event := &example.Example{Message: "Hello Word"}

	if err := producer.Publish(ctx, event); err != nil {
		return
	}

	slog.Info("event sent")
}
