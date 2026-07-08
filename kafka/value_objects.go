package kafka

import (
	"github.com/lcnascimento/go-kit/env"
	"github.com/segmentio/kafka-go"
)

var (
	brokers = env.GetList("KAFKA_BROKERS", env.WithDefaultListValue([]string{"localhost:9092"}))
	groupID = env.Get("OTEL_SERVICE_NAME", env.WithDefaultValue("default"))
)

type (
	Message      = kafka.Message
	EventType    string
	EventVersion string
)

type Event interface {
	GetTopic() string
	GetKey() []byte
	GetType() EventType
	GetVersion() EventVersion
}
