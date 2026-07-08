package kafka

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/kafka"
	logger = log.MustNewLogger(pkg)
	tracer = otel.Tracer(pkg)
)

type debugLogger struct {
	logger log.Logger
}

func newDebugLogger() *debugLogger {
	return &debugLogger{logger: logger}
}

func (l *debugLogger) Printf(msg string, args ...interface{}) {
	l.logger.Debug(context.Background(), fmt.Sprintf(msg, args...))
}

type errorLogger struct {
	logger log.Logger
}

func newErrorLogger() *errorLogger {
	return &errorLogger{logger: logger}
}

func (l *errorLogger) Printf(msg string, args ...interface{}) {
	l.logger.ErrorMessage(context.Background(), fmt.Sprintf(msg, args...))
}

func (s *Producer) onPublishStart(ctx context.Context, count int) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "PublishMessages")
	logger.Debug(ctx, "publishing messages", log.Int("count", count))

	return ctx, span
}

func (s *Producer) onStop(ctx context.Context) {
	logger.Info(ctx, "stopping kafka producer")
}

func (s *Subscriber[T]) onStart(ctx context.Context, topic string) {
	logger.Info(ctx, "starting kafka subscriber", log.String("topic", topic))
}

func (s *Subscriber[T]) onStop(ctx context.Context, topic string) {
	logger.Info(ctx, "stopping kafka subscriber", log.String("topic", topic))
}

func (s *Subscriber[T]) onConsumeStart(ctx context.Context, msg Message) (context.Context, trace.Span) {
	carrier := propagation.MapCarrier{}
	for _, h := range msg.Headers {
		carrier.Set(h.Key, string(h.Value))
	}

	propagator := otel.GetTextMapPropagator()
	wireCtx := propagator.Extract(context.Background(), carrier)

	ctx, span := tracer.Start(ctx, "ConsumeMessage", trace.WithLinks(trace.LinkFromContext(wireCtx)))
	logger.Debug(
		ctx, "consuming kafka message",
		log.String("message.topic", msg.Topic),
		log.String("message.key", string(msg.Key)),
		log.Int("message.partition", msg.Partition),
	)

	return ctx, span
}

func (s *Subscriber[T]) onError(ctx context.Context, err error) error {
	logger.ErrorBySeverity(ctx, err)

	return err
}

func (s *Subscriber[T]) onErrorWithSpan(ctx context.Context, err error, span trace.Span) error {
	span.End()

	return s.onError(ctx, err)
}
