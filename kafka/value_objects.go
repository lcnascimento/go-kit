package kafka

import (
	"github.com/lcnascimento/go-kit/env"
)

var (
	brokers = env.GetList("KAFKA_BROKERS", env.WithDefaultListValue([]string{"localhost:9092"}))
	groupID = env.Get("OTEL_SERVICE_NAME", env.WithDefaultValue("default"))
)

type (
	EventType    string
	EventVersion string
)

func (t EventType) String() string {
	return string(t)
}

func (t EventVersion) String() string {
	return string(t)
}

type Event interface {
	GetTopic() string
	GetKey() []byte
	GetType() EventType
	GetVersion() EventVersion
}
