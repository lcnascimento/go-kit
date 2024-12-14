package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/lcnascimento/go-kit/messaging"
	"github.com/lcnascimento/go-kit/o11y"
)

type broker struct {
	logger           watermill.LoggerAdapter
	pubsub           messaging.PubSub
	router           *message.Router
	marshaler        cqrs.CommandEventMarshaler
	commandProcessor *cqrs.CommandProcessor
	commandBus       *cqrs.CommandBus
	eventProcessor   *cqrs.EventProcessor
	eventBus         *cqrs.EventBus
}

// NewBroker creates a new BrokerCQRS instance.
func NewBroker(opts ...Option) (messaging.BrokerCQRS, error) {
	var err error

	ctx := o11y.Context()

	broker := &broker{
		logger:    messaging.NewWatermillLogger(),
		marshaler: cqrs.JSONMarshaler{},
	}

	for _, opt := range opts {
		opt(broker)
	}

	broker.pubsub = gochannel.NewGoChannel(
		gochannel.Config{},
		broker.logger,
	)

	broker.router, err = newRouterBuilder(broker.logger).build(ctx)
	if err != nil {
		return nil, err
	}

	broker.commandBus, err = newCommandBusBuilder(
		broker.pubsub,
		broker.marshaler,
		broker.logger,
	).build(ctx)
	if err != nil {
		return nil, err
	}

	broker.commandProcessor, err = newCommandProcessorBuilder(
		broker.router,
		broker.pubsub,
		broker.marshaler,
		broker.logger,
	).build(ctx)
	if err != nil {
		return nil, err
	}

	broker.eventBus, err = newEventBusBuilder(
		broker.pubsub,
		broker.marshaler,
		broker.logger,
	).build(ctx)
	if err != nil {
		return nil, err
	}

	broker.eventProcessor, err = newEventProcessorBuilder(
		broker.pubsub,
		broker.router,
		broker.marshaler,
		broker.logger,
	).build(ctx)
	if err != nil {
		return nil, err
	}

	return broker, nil
}

// Start starts the broker.
// All Command and Event handlers must be added before calling this method.
func (b *broker) Start(ctx context.Context) (err error) {
	b.onStart(ctx)
	return b.router.Run(ctx)
}

// Stop stops the broker.
func (b *broker) Stop(ctx context.Context) error {
	b.onStop(ctx)
	return b.router.Close()
}

// Running is closed when broker is running.
// In other words: you can wait till broker is running using
//
//	fmt.Println("Starting broker")
//	go r.Run(ctx)
//	<- r.Running()
//	fmt.Println("Broker is running")
//
// Warning: for historical reasons, this channel is not aware of broker closing.
// The channel will be closed if the broker has been running and closed.
func (b *broker) Running(ctx context.Context) chan struct{} {
	out := make(chan struct{})
	go func() {
		<-b.router.Running()
		b.onRunning(ctx)
		close(out)
	}()

	return out
}

// AddCommandHandlers adds command handlers to the command processor.
func (b *broker) AddCommandHandlers(ctx context.Context, handlers ...cqrs.CommandHandler) error {
	if err := b.commandProcessor.AddHandlers(handlers...); err != nil {
		return b.onAddCommandHandlersError(ctx, err)
	}

	return nil
}

// AddEventHandlers adds event handlers to the event processor.
func (b *broker) AddEventHandlers(ctx context.Context, handlers ...cqrs.EventHandler) error {
	if err := b.eventProcessor.AddHandlers(handlers...); err != nil {
		return b.onAddEventHandlersError(ctx, err)
	}

	return nil
}

// SendCommand sends a command to the command bus.
func (b *broker) SendCommand(ctx context.Context, command any) error {
	if err := b.commandBus.Send(ctx, command); err != nil {
		return b.onSendCommandError(ctx, command, err)
	}

	return nil
}

// SendEvent sends an event to the event bus.
func (b *broker) SendEvent(ctx context.Context, event any) error {
	if err := b.eventBus.Publish(ctx, event); err != nil {
		return b.onSendEventError(ctx, event, err)
	}

	return nil
}
