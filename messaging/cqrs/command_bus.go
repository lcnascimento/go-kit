package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"github.com/lcnascimento/go-kit/messaging"
)

type commandBusBuilder struct {
	pubsub    messaging.PubSub
	marshaler cqrs.CommandEventMarshaler
	logger    watermill.LoggerAdapter
}

func newCommandBusBuilder(
	pubsub messaging.PubSub,
	marshaler cqrs.CommandEventMarshaler,
	logger watermill.LoggerAdapter,
) *commandBusBuilder {
	return &commandBusBuilder{
		pubsub:    pubsub,
		marshaler: marshaler,
		logger:    logger,
	}
}

func (b *commandBusBuilder) build(ctx context.Context) (*cqrs.CommandBus, error) {
	commandBus, err := cqrs.NewCommandBusWithConfig(b.pubsub, cqrs.CommandBusConfig{
		GeneratePublishTopic: b.generateCommandBusPublishTopic,
		Marshaler:            b.marshaler,
		Logger:               b.logger,
	})
	if err != nil {
		return nil, b.onBuildError(ctx, err)
	}

	return commandBus, nil
}

func (b *commandBusBuilder) generateCommandBusPublishTopic(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
	return params.CommandName, nil
}
