package messaging

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

// PubSub is a contract for a message broker that implements both Publisher and Subscriber.
type PubSub interface {
	message.Publisher
	message.Subscriber
}

// BrokerCQRS is a contract for a message broker that implements the CQRS pattern.
type BrokerCQRS interface {
	// Start starts the broker.
	// All Command and Event handlers must be added before calling this method.
	Start(ctx context.Context) error

	// Stop stops the broker.
	Stop(ctx context.Context) error

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
	Running(ctx context.Context) chan struct{}

	// AddCommandHandlers adds command handlers to the command processor.
	AddCommandHandlers(handlers ...cqrs.CommandHandler)

	// AddEventHandlers adds event handlers to the event processor.
	AddEventHandlers(handlers ...cqrs.EventHandler)

	// SendCommand sends a command to the command bus.
	SendCommand(ctx context.Context, command any) error

	// SendEvent sends an event to the event bus.
	SendEvent(ctx context.Context, event any) error
}
