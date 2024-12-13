package messaging

import "github.com/ThreeDotsLabs/watermill"

const (
	logKeyCommand = "command"
	logKeyEvent   = "event"
)

// Option is a function that configures the BrokerCQRS.
type Option func(broker *brokerCQRS)

// WithLogger configures the logger.
func WithLogger(logger watermill.LoggerAdapter) Option {
	return func(broker *brokerCQRS) {
		broker.logger = logger
	}
}
