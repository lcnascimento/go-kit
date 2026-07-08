package kafka

import "github.com/segmentio/kafka-go"

type messageCarrier struct {
	msg *kafka.Message
}

func (c messageCarrier) Get(key string) string {
	for _, h := range c.msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}

	return ""
}

func (c messageCarrier) Set(key string, value string) {
	for i, h := range c.msg.Headers {
		if h.Key == key {
			c.msg.Headers[i].Value = []byte(value)
			return
		}
	}

	c.msg.Headers = append(c.msg.Headers, kafka.Header{Key: key, Value: []byte(value)})
}

func (c messageCarrier) Keys() []string {
	keys := make([]string, len(c.msg.Headers))
	for i, h := range c.msg.Headers {
		keys[i] = h.Key
	}

	return keys
}
