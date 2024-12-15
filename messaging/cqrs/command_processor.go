package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/codes"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/messaging"
)

type commandProcessorBuilder struct {
	router    *message.Router
	pubsub    messaging.PubSub
	marshaler cqrs.CommandEventMarshaler
	logger    watermill.LoggerAdapter
}

func newCommandProcessorBuilder(
	router *message.Router,
	pubsub messaging.PubSub,
	marshaler cqrs.CommandEventMarshaler,
	logger watermill.LoggerAdapter,
) *commandProcessorBuilder {
	return &commandProcessorBuilder{
		router:    router,
		pubsub:    pubsub,
		marshaler: marshaler,
		logger:    logger,
	}
}

func (b *commandProcessorBuilder) build(ctx context.Context) (*cqrs.CommandProcessor, error) {
	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(
		b.router,
		cqrs.CommandProcessorConfig{
			GenerateSubscribeTopic: b.generateCommandProcessorSubscribeTopic,
			SubscriberConstructor:  b.commandProcessorSubscriberConstructor,
			OnHandle:               b.onHandle,
			Marshaler:              b.marshaler,
			Logger:                 b.logger,
		},
	)
	if err != nil {
		return nil, b.onBuildError(ctx, err)
	}

	return commandProcessor, nil
}

func (b *commandProcessorBuilder) generateCommandProcessorSubscribeTopic(
	params cqrs.CommandProcessorGenerateSubscribeTopicParams,
) (string, error) {
	return params.CommandName, nil
}

func (b *commandProcessorBuilder) commandProcessorSubscriberConstructor(
	_ cqrs.CommandProcessorSubscriberConstructorParams,
) (message.Subscriber, error) {
	return b.pubsub, nil
}

func (b *commandProcessorBuilder) onHandle(params cqrs.CommandProcessorOnHandleParams) (err error) {
	ctx, span := b.onHandleStart(params)
	defer b.onHandleEnd(ctx, span, params)

	err = params.Handler.Handle(ctx, params.Command)
	if err == nil {
		span.SetStatus(codes.Ok, "command handled")
		return nil
	}

	if errors.IsRetryable(err) {
		return err
	}

	return nil
}
