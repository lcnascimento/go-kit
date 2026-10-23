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

// WithShutdown sets shutdown, called when the forwarder stops without Stop being called.
func WithShutdown(shutdown func()) ForwarderOption {
	return func(f *Forwarder) {
		f.shutdown = shutdown
	}
}

// WithValidation makes the forwarder send only messages validate accepts; it retries a rejected one, as it does a
// failed send, so a malformed message holds back the ones after it instead of reaching Kafka.
func WithValidation(validate func(*message.Message) error) ForwarderOption {
	return func(f *Forwarder) {
		f.validate = validate
	}
}
