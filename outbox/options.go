package outbox

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
)

type JanitorOption func(*Janitor)

// WithRetention keeps forwarded messages for retention before purging them; the default is 7 days.
func WithRetention(retention time.Duration) JanitorOption {
	return func(j *Janitor) {
		j.retention = retention
	}
}

// WithPurgeInterval sets how often the Janitor purges; the default is 1 hour.
func WithPurgeInterval(interval time.Duration) JanitorOption {
	return func(j *Janitor) {
		j.interval = interval
	}
}

type ForwarderOption func(*Forwarder)

// WithShutdown sets the function called when the forwarder stops without Stop.
func WithShutdown(shutdown func()) ForwarderOption {
	return func(f *Forwarder) {
		f.shutdown = shutdown
	}
}

// WithValidation retries a message validate rejects, as it does a failed send, so it holds back the messages after it.
func WithValidation(validate func(*message.Message) error) ForwarderOption {
	return func(f *Forwarder) {
		f.validate = validate
	}
}
