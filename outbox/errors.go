package outbox

import "github.com/lcnascimento/go-kit/errors"

var (
	ErrInvalidName = errors.New("the outbox name must be a lowercase identifier: it names its tables and metrics").
			WithCode("INVALID_OUTBOX_NAME")

	ErrPublish = errors.New("failed to write the message to the outbox").WithCode("OUTBOX_PUBLISH")

	ErrPurge = errors.New("failed to purge the outbox").WithCode("OUTBOX_PURGE")

	ErrOldestPendingAge = errors.New("failed to read the age of the oldest pending outbox message").
				WithCode("OUTBOX_OLDEST_PENDING_AGE")

	ErrStartJanitor = errors.New("outbox janitor did not start").WithCode("OUTBOX_JANITOR_START")

	ErrKafkaUnavailable = errors.New("failed to connect to Kafka: are the brokers up and their addresses right?").
				WithCode("KAFKA_UNAVAILABLE")

	ErrTopicMissing = errors.New("topic does not exist: create it before starting the forwarder").
			WithCode("TOPIC_MISSING")

	ErrInvalidForwarderConfig = errors.New("forwarder needs the outbox pool, Kafka brokers and the topics it forwards to").
					WithCode("INVALID_FORWARDER_CONFIG")

	ErrStartForwarder = errors.New("forwarder did not start").WithCode("FORWARDER_START")

	ErrForwarderNotRunning = errors.New("forwarder is not running").WithCode("FORWARDER_NOT_RUNNING")

	ErrStopForwarder = errors.New("forwarder did not stop cleanly").WithCode("FORWARDER_STOP")

	ErrSend = errors.New("failed to send the message to Kafka: retrying").WithCode("FORWARDER_SEND")
)
