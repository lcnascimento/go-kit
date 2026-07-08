package example

import "github.com/lcnascimento/go-kit/kafka"

type Example struct {
	Message string `json:"message"`
}

func (e *Example) GetTopic() string {
	return "example"
}

func (e *Example) GetKey() []byte {
	return []byte{}
}

func (e *Example) GetType() kafka.EventType {
	return kafka.EventType("EXAMPLE")
}

func (e *Example) GetVersion() kafka.EventVersion {
	return kafka.EventVersion("V1")
}
