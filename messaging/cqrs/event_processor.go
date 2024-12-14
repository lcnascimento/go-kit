package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/messaging"
	"go.opentelemetry.io/otel/codes"
)

type eventProcessorBuilder struct {
	pubsub    messaging.PubSub
	router    *message.Router
	marshaler cqrs.CommandEventMarshaler
	logger    watermill.LoggerAdapter
}

func newEventProcessorBuilder(
	pubsub messaging.PubSub,
	router *message.Router,
	marshaler cqrs.CommandEventMarshaler,
	logger watermill.LoggerAdapter,
) *eventProcessorBuilder {
	return &eventProcessorBuilder{
		pubsub:    pubsub,
		router:    router,
		marshaler: marshaler,
		logger:    logger,
	}
}

func (b *eventProcessorBuilder) build(ctx context.Context) (*cqrs.EventProcessor, error) {
	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		b.router,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: b.generateSubscribeTopic,
			SubscriberConstructor:  b.subscriberConstructor,
			OnHandle:               b.onHandle,
			Marshaler:              b.marshaler,
			Logger:                 b.logger,
		},
	)
	if err != nil {
		return nil, b.onBuildError(ctx, err)
	}

	return eventProcessor, nil
}

func (b *eventProcessorBuilder) generateSubscribeTopic(
	params cqrs.EventProcessorGenerateSubscribeTopicParams,
) (string, error) {
	return params.EventName, nil
}

func (b *eventProcessorBuilder) subscriberConstructor(
	params cqrs.EventProcessorSubscriberConstructorParams,
) (message.Subscriber, error) {
	return b.pubsub, nil
}

func (b *eventProcessorBuilder) onHandle(params cqrs.EventProcessorOnHandleParams) (err error) {
	ctx, span := b.onStart(params)
	defer b.onEnd(ctx, span, params)

	err = params.Handler.Handle(ctx, params.Event)
	if err == nil {
		span.SetStatus(codes.Ok, "event handled")
		return nil
	}

	if errors.IsRetryable(err) {
		return err
	}

	return nil
}
