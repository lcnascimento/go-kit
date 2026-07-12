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
			StartOffset: kafka.FirstOffset,
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

const (
	defaultMaxBatchSize = 100
	defaultMaxBatchWait = 5 * time.Second
)

type batchOptions struct {
	maxSize int
	maxWait time.Duration
}

type BatchOption func(*batchOptions)

func WithMaxBatchSize(size int) BatchOption {
	return func(o *batchOptions) {
		if size > 0 {
			o.maxSize = size
		}
	}
}

func WithMaxBatchWait(wait time.Duration) BatchOption {
	return func(o *batchOptions) {
		if wait > 0 {
			o.maxWait = wait
		}
	}
}

func (s *Subscriber[T]) RunInBatch(ctx context.Context, cb func(context.Context, []T) error, opts ...BatchOption) error {
	options := &batchOptions{
		maxSize: defaultMaxBatchSize,
		maxWait: defaultMaxBatchWait,
	}
	for _, opt := range opts {
		opt(options)
	}

	s.onStart(ctx, s.topic)

	for {
		msgs, values, err := s.fetchBatch(ctx, options)
		if err != nil {
			return s.onError(ctx, err)
		}

		if ctx.Err() != nil {
			return nil
		}

		if len(msgs) == 0 {
			continue
		}

		ctx, span := s.onConsumeBatchStart(ctx, msgs)

		if err := cb(ctx, values); err != nil {
			span.End()
			return err
		}

		err = s.reader.CommitMessages(ctx, msgs...)
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

func (s *Subscriber[T]) fetchBatch(ctx context.Context, options *batchOptions) ([]kafka.Message, []T, error) {
	batchCtx, cancel := context.WithTimeout(ctx, options.maxWait)
	defer cancel()

	msgs := make([]kafka.Message, 0, options.maxSize)
	values := make([]T, 0, options.maxSize)

	for len(msgs) < options.maxSize {
		msg, err := s.reader.FetchMessage(batchCtx)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return msgs, values, nil
		}
		if err != nil {
			return nil, nil, err
		}

		var value T
		if err := json.Unmarshal(msg.Value, &value); err != nil {
			return nil, nil, err
		}

		msgs = append(msgs, msg)
		values = append(values, value)
	}

	return msgs, values, nil
}

func (s *Subscriber[T]) Stop(ctx context.Context) error {
	s.onStop(ctx, s.topic)

	return s.reader.Close()
}
