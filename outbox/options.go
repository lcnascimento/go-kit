package outbox

import "time"

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

// WithShutdown is called when the forwarder stops without Stop being called.
func WithShutdown(shutdown func()) ForwarderOption {
	return func(f *Forwarder) {
		f.shutdown = shutdown
	}
}
