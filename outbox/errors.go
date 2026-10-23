package outbox

import "github.com/lcnascimento/go-kit/errors"

var (
	ErrInvalidName = errors.New("the outbox name must be a lowercase identifier of at most 38 characters").
			WithKind(errors.KindInvalidInput).
			WithCode("INVALID_OUTBOX_NAME")

	ErrInvalidForwarderConfig = errors.New("the forwarder needs the outbox pool, Kafka brokers and the topics it forwards to").
					WithKind(errors.KindInvalidInput).
					WithCode("INVALID_FORWARDER_CONFIG")

	ErrInvalidMessage = errors.New("the outbox message failed the forwarder's validation").
				WithCode("INVALID_OUTBOX_MESSAGE")

	ErrKafkaUnavailable = errors.New("failed to reach Kafka: are the brokers up and their addresses right?").
				WithKind(errors.KindServiceUnavailable).
				WithCode("KAFKA_UNAVAILABLE").
				Retryable()

	ErrTopicMissing = errors.New("topic does not exist: create it before starting the forwarder").
			WithCode("TOPIC_MISSING")

	ErrForwarderNotRunning = errors.New("forwarder is not running").WithCode("FORWARDER_NOT_RUNNING")

	ErrAlreadyStarted = errors.New("already started: create a new one to start again").WithCode("ALREADY_STARTED")

	ErrJanitorAlreadyRunning = errors.New("another janitor of this outbox is already running in the process").
					WithCode("JANITOR_ALREADY_RUNNING")

	ErrPublish          = errors.New("failed to publish the message to the outbox")
	ErrPurge            = errors.New("failed to purge the outbox")
	ErrOldestPendingAge = errors.New("failed to read the oldest pending age of the outbox")
	ErrNewJanitor       = errors.New("failed to create the janitor")
	ErrStartJanitor     = errors.New("failed to start the janitor")
	ErrStopJanitor      = errors.New("failed to stop the janitor")
	ErrStartForwarder   = errors.New("failed to start the forwarder")
	ErrStopForwarder    = errors.New("failed to stop the forwarder")
	ErrSend             = errors.New("failed to send the message to Kafka").Retryable()
)
