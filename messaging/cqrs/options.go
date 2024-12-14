package cqrs

import "github.com/ThreeDotsLabs/watermill"

// Option is a function that configures the BrokerCQRS.
type Option func(broker *broker)

// WithLogger configures the logger.
func WithLogger(logger watermill.LoggerAdapter) Option {
	return func(broker *broker) {
		broker.logger = logger
	}
}
