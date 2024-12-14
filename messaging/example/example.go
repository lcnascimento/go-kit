package main

import (
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"

	"github.com/lcnascimento/go-kit/messaging/cqrs"
	"github.com/lcnascimento/go-kit/messaging/example/command"
	"github.com/lcnascimento/go-kit/messaging/example/event"
	"github.com/lcnascimento/go-kit/o11y"
	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/messaging/example"
	logger = log.NewLogger(pkg)
	tracer = otel.Tracer(pkg)
)

func main() {
	defer o11y.Shutdown()
	ctx := o11y.Context()

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	foo, _ := baggage.NewMember("foo", "foo")
	bar, _ := baggage.NewMember("bar", "bar")

	bag, _ := baggage.New(foo, bar)

	ctx = baggage.ContextWithBaggage(ctx, bag)

	broker, err := cqrs.NewBroker()
	if err != nil {
		return
	}

	go func() {
		broker.AddCommandHandlers(ctx, &command.CommandLogMessageHandler{Broker: broker, Logger: logger})
		broker.AddEventHandlers(ctx, &event.EventMessageLoggedHandler{Logger: logger})
		broker.Start(ctx)
	}()

	<-broker.Running(ctx)

	if err := broker.SendCommand(ctx, &command.CommandLogMessage{Message: "Hello, World!"}); err != nil {
		return
	}
	logger.Info(ctx, "command sent")

	time.Sleep(1 * time.Second)

	broker.Stop(ctx)
}
