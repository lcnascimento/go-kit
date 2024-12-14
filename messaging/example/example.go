package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"

	"github.com/lcnascimento/go-kit/messaging"
	"github.com/lcnascimento/go-kit/o11y"
	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/messaging/example"
	logger = log.NewLogger(pkg)
	tracer = otel.Tracer(pkg)
)

type event struct {
	Message string `json:"message"`
}

type eventHandler struct{}

func (h *eventHandler) HandlerName() string {
	return "example"
}

func (h *eventHandler) NewEvent() any {
	return &event{}
}

func (h *eventHandler) Handle(ctx context.Context, e any) error {
	event, _ := e.(*event)
	logger.Info(ctx, event.Message)
	return nil
}

func main() {
	defer o11y.Shutdown()
	ctx := o11y.Context()

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	foo, _ := baggage.NewMember("foo", "foo")
	bar, _ := baggage.NewMember("bar", "bar")

	bag, _ := baggage.New(foo, bar)

	ctx = baggage.ContextWithBaggage(ctx, bag)

	broker, err := messaging.NewBrokerCQRS()
	if err != nil {
		return
	}

	go func() {
		broker.AddEventHandlers(ctx, &eventHandler{})
		broker.Start(ctx)
	}()

	<-broker.Running(ctx)

	if err := broker.SendEvent(ctx, &event{Message: "Hello, World!"}); err != nil {
		return
	}
	logger.Info(ctx, "event sent")

	broker.Stop(ctx)
}
