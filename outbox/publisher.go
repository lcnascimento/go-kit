package outbox

import (
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
)

type tracedPublisher struct {
	message.Publisher

	outboxName string
	validate   func(*message.Message) error

	mu       sync.Mutex
	lastID   string
	attempts int
}

func newTracedPublisher(publisher message.Publisher, outboxName string, validate func(*message.Message) error) *tracedPublisher {
	return &tracedPublisher{Publisher: publisher, outboxName: outboxName, validate: validate}
}

func (p *tracedPublisher) Publish(topic string, messages ...*message.Message) error {
	for _, msg := range messages {
		end := onSendStart(msg, p.outboxName, topic, p.attempt(msg.UUID))

		err := p.send(topic, msg)
		end(err)

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *tracedPublisher) send(topic string, msg *message.Message) error {
	if p.validate != nil {
		if err := p.validate(msg); err != nil {
			return ErrInvalidMessage.WithCause(err)
		}
	}

	return p.Publisher.Publish(topic, msg)
}

func (p *tracedPublisher) attempt(messageID string) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	if messageID != p.lastID {
		p.lastID = messageID
		p.attempts = 0
	}

	p.attempts++

	return p.attempts
}
