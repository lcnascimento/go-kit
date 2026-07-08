package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/lcnascimento/go-kit/errors"
)

type Subscriber[T Event] struct {
	topic  string
	reader kafka.Reader
}

func NewSubscriber[T Event](topic string) *Subscriber[T] {
	const (
		defaultReadTimeout = 10 * time.Second
		defaultDialTimeout = 10 * time.Second
	)

	subs := &Subscriber[T]{
		topic: topic,
		reader: *kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			GroupID:     groupID,
			Topic:       topic,
			StartOffset: kafka.LastOffset,
			MaxWait:     defaultReadTimeout,
			Logger:      newDebugLogger(),
			ErrorLogger: newErrorLogger(),
			Dialer:      &kafka.Dialer{Timeout: defaultDialTimeout},
		}),
	}

	return subs
}

func (s *Subscriber[T]) Run(ctx context.Context, cb func(context.Context, T) error) error {
	s.onStart(ctx, s.topic)

	for {
		msg, err := s.reader.FetchMessage(ctx)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil
		}
		if err != nil {
			return s.onError(ctx, err)
		}

		var value T
		if err := json.Unmarshal(msg.Value, &value); err != nil {
			return s.onError(ctx, err)
		}

		ctx, span := s.onConsumeStart(ctx, value.GetType(), msg)

		if err := cb(ctx, value); err != nil {
			span.End()
			return err
		}

		err = s.reader.CommitMessages(ctx, msg)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			span.End()
			return nil
		}
		if err != nil {
			return s.onErrorWithSpan(ctx, err, span)
		}

		span.End()
	}
}

func (s *Subscriber[T]) Stop(ctx context.Context) error {
	s.onStop(ctx, s.topic)

	return s.reader.Close()
}
