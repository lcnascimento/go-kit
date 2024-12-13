package messaging

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/lcnascimento/go-kit/o11y"
)

type brokerCQRS struct {
	logger           watermill.LoggerAdapter
	pubsub           PubSub
	router           *message.Router
	marshaler        cqrs.CommandEventMarshaler
	commandProcessor *cqrs.CommandProcessor
	commandBus       *cqrs.CommandBus
	eventProcessor   *cqrs.EventProcessor
	eventBus         *cqrs.EventBus
}

// NewBrokerCQRS creates a new BrokerCQRS instance.
func NewBrokerCQRS(opts ...Option) (BrokerCQRS, error) {
	ctx := o11y.Context()

	broker := &brokerCQRS{
		logger:    watermill.NewSlogLogger(nil),
		marshaler: cqrs.JSONMarshaler{},
	}

	for _, opt := range opts {
		opt(broker)
	}

	broker.pubsub = gochannel.NewGoChannel(
		gochannel.Config{},
		broker.logger,
	)

	if err := broker.buildRouter(ctx); err != nil {
		return nil, err
	}

	if err := broker.buildCommandBus(ctx); err != nil {
		return nil, err
	}

	if err := broker.buildCommandProcessor(ctx); err != nil {
		return nil, err
	}

	if err := broker.buildEventBus(ctx); err != nil {
		return nil, err
	}

	if err := broker.buildEventProcessor(ctx); err != nil {
		return nil, err
	}

	return broker, nil
}

// Start starts the broker.
// All Command and Event handlers must be added before calling this method.
func (b *brokerCQRS) Start(ctx context.Context) (err error) {
	b.onStart(ctx)
	return b.router.Run(ctx)
}

// Stop stops the broker.
func (b *brokerCQRS) Stop(ctx context.Context) error {
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
func (b *brokerCQRS) Running(ctx context.Context) chan struct{} {
	out := make(chan struct{})
	go func() {
		<-b.router.Running()
		b.onRunning(ctx)
		close(out)
	}()

	return out
}

// AddCommandHandlers adds command handlers to the command processor.
func (b *brokerCQRS) AddCommandHandlers(ctx context.Context, handlers ...cqrs.CommandHandler) error {
	if err := b.commandProcessor.AddHandlers(handlers...); err != nil {
		return b.onAddCommandHandlersError(ctx, err)
	}

	return nil
}

// AddEventHandlers adds event handlers to the event processor.
func (b *brokerCQRS) AddEventHandlers(ctx context.Context, handlers ...cqrs.EventHandler) error {
	if err := b.eventProcessor.AddHandlers(handlers...); err != nil {
		return b.onAddEventHandlersError(ctx, err)
	}

	return nil
}

// SendCommand sends a command to the command bus.
func (b *brokerCQRS) SendCommand(ctx context.Context, command any) error {
	if err := b.commandBus.Send(ctx, command); err != nil {
		return b.onSendCommandError(ctx, command, err)
	}

	return nil
}

// SendEvent sends an event to the event bus.
func (b *brokerCQRS) SendEvent(ctx context.Context, event any) error {
	if err := b.eventBus.Publish(ctx, event); err != nil {
		return b.onSendEventError(ctx, event, err)
	}

	return nil
}

func (b *brokerCQRS) buildRouter(ctx context.Context) error {
	var err error

	b.router, err = message.NewRouter(message.RouterConfig{}, b.logger)
	if err != nil {
		return b.onBuildRouterError(ctx, err)
	}
	b.router.AddMiddleware(middleware.Recoverer)

	return nil
}

func (b *brokerCQRS) buildCommandBus(ctx context.Context) error {
	var err error

	b.commandBus, err = cqrs.NewCommandBusWithConfig(b.pubsub, cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return params.CommandName, nil
		},
		Marshaler: b.marshaler,
	})
	if err != nil {
		return b.onBuildCommandBusError(ctx, err)
	}

	return nil
}

func (b *brokerCQRS) buildCommandProcessor(ctx context.Context) error {
	var err error

	b.commandProcessor, err = cqrs.NewCommandProcessorWithConfig(
		b.router,
		cqrs.CommandProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.CommandName, nil
			},
			SubscriberConstructor: func(_ cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return b.pubsub, nil
			},
			Marshaler: b.marshaler,
		},
	)
	if err != nil {
		return b.onBuildCommandProcessorError(ctx, err)
	}

	return nil
}

func (b *brokerCQRS) buildEventBus(ctx context.Context) error {
	var err error

	b.eventBus, err = cqrs.NewEventBusWithConfig(b.pubsub, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return params.EventName, nil
		},
		Marshaler: b.marshaler,
	})
	if err != nil {
		return b.onBuildEventBusError(ctx, err)
	}

	return nil
}

func (b *brokerCQRS) buildEventProcessor(ctx context.Context) error {
	var err error

	b.eventProcessor, err = cqrs.NewEventProcessorWithConfig(
		b.router,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			SubscriberConstructor: func(_ cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return b.pubsub, nil
			},
			Marshaler: b.marshaler,
		},
	)
	if err != nil {
		return b.onBuildEventProcessorError(ctx, err)
	}

	return nil
}
