package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/lcnascimento/go-kit/messaging"
)

type commandBusBuilder struct {
	pubsub     messaging.PubSub
	marshaler  cqrs.CommandEventMarshaler
	logger     watermill.LoggerAdapter
	propagator propagation.TextMapPropagator
}

func newCommandBusBuilder(
	pubsub messaging.PubSub,
	marshaler cqrs.CommandEventMarshaler,
	logger watermill.LoggerAdapter,
) *commandBusBuilder {
	return &commandBusBuilder{
		pubsub:     pubsub,
		marshaler:  marshaler,
		logger:     logger,
		propagator: otel.GetTextMapPropagator(),
	}
}

func (b *commandBusBuilder) build(ctx context.Context) (*cqrs.CommandBus, error) {
	commandBus, err := cqrs.NewCommandBusWithConfig(b.pubsub, cqrs.CommandBusConfig{
		GeneratePublishTopic: b.generateCommandBusPublishTopic,
		OnSend:               b.onSend,
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

func (b *commandBusBuilder) onSend(params cqrs.CommandBusOnSendParams) error {
	ctx, span := b.onSendStart(params)
	defer b.onSendEnd(ctx, span)

	b.propagator.Inject(ctx, propagation.MapCarrier(params.Message.Metadata))
	params.Message.SetContext(ctx)

	return nil
}
