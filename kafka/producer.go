package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	const defaultWriteTimeout = 10 * time.Second

	producer := &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: defaultWriteTimeout,
			Logger:       newDebugLogger(),
			ErrorLogger:  newErrorLogger(),
		},
	}

	return producer
}

func (p *Producer) Publish(ctx context.Context, events ...Event) error {
	ctx, span := p.onPublishStart(ctx, len(events))
	defer span.End()

	messages := make([]kafka.Message, 0, len(events))

	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return errors.ErrCastPayload.WithCause(err)
		}

		msg := kafka.Message{
			Topic: event.GetTopic(),
			Key:   event.GetKey(),
			Value: payload,
		}

		otel.GetTextMapPropagator().Inject(ctx, &messageCarrier{msg: &msg})
		messages = append(messages, msg)
	}

	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		return ErrWriteMessages.WithCause(err)
	}

	return nil
}

func (p *Producer) Stop(ctx context.Context) error {
	p.onStop(ctx)

	return p.writer.Close()
}
