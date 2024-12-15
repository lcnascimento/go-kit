package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/lcnascimento/go-kit/messaging"
)

type eventBusBuilder struct {
	pubsub     messaging.PubSub
	marshaler  cqrs.CommandEventMarshaler
	logger     watermill.LoggerAdapter
	propagator propagation.TextMapPropagator
}

func newEventBusBuilder(
	pubsub messaging.PubSub,
	marshaler cqrs.CommandEventMarshaler,
	logger watermill.LoggerAdapter,
) *eventBusBuilder {
	return &eventBusBuilder{
		pubsub:     pubsub,
		marshaler:  marshaler,
		logger:     logger,
		propagator: otel.GetTextMapPropagator(),
	}
}

func (b *eventBusBuilder) build(ctx context.Context) (*cqrs.EventBus, error) {
	eventBus, err := cqrs.NewEventBusWithConfig(b.pubsub, cqrs.EventBusConfig{
		GeneratePublishTopic: b.generatePublishTopic,
		OnPublish:            b.onPublish,
		Marshaler:            b.marshaler,
		Logger:               b.logger,
	})
	if err != nil {
		return nil, b.onBuildError(ctx, err)
	}

	return eventBus, nil
}

func (b *eventBusBuilder) generatePublishTopic(params cqrs.GenerateEventPublishTopicParams) (string, error) {
	return params.EventName, nil
}

func (b *eventBusBuilder) onPublish(params cqrs.OnEventSendParams) error {
	ctx, span := b.onPublishStart(params)
	defer b.onPublishEnd(ctx, span)

	b.propagator.Inject(ctx, propagation.MapCarrier(params.Message.Metadata))
	params.Message.SetContext(ctx)

	return nil
}
