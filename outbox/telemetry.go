package outbox

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	otelmetric "go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/log"
	"github.com/lcnascimento/go-kit/o11y/metric"
)

const (
	operationCreate = "create"
	operationSend   = "send"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/outbox"
	logger = log.MustNewLogger(pkg)
	meter  = otel.Meter(pkg)
	tracer = otel.Tracer(pkg)

	sentMessages = metric.MustIntCounter(meter, "messaging.client.sent.messages",
		"Number of attempts to send a message to Kafka")
	operationDuration = metric.MustDurationHistogram(meter, "messaging.client.operation.duration",
		"Duration of one attempt to send a message to Kafka")

	outcomeSuccess = otelmetric.WithAttributes(attribute.String("outcome", "success"))
	outcomeFailure = otelmetric.WithAttributes(attribute.String("outcome", "failure"))
)

type janitorTelemetry struct {
	purgedMessages   otelmetric.Int64Counter
	purgeRuns        otelmetric.Int64Counter
	oldestPendingAge otelmetric.Int64ObservableGauge
}

// Exception to package-level instruments named <resource>.<event>: the Janitor's are per outbox and named
// <name>.outbox.*, because the boilerplate's alerts, SLOs and dashboards query <service>_outbox_* per service
// (OutboxAgeAbsent among them) and must keep working for the services already on it.
func newJanitorTelemetry(name string) (janitorTelemetry, error) {
	purgedMessages, err := meter.Int64Counter(name+".outbox.purged.messages",
		otelmetric.WithDescription("Number of forwarded messages deleted from the outbox"))
	if err != nil {
		return janitorTelemetry{}, err
	}

	purgeRuns, err := meter.Int64Counter(name+".outbox.purge.runs",
		otelmetric.WithDescription("Number of outbox purges, by outcome"))
	if err != nil {
		return janitorTelemetry{}, err
	}

	oldestPendingAge, err := meter.Int64ObservableGauge(name+".outbox.oldest_pending.age",
		otelmetric.WithDescription("Age of the oldest outbox message not yet forwarded to Kafka, 0 when none"),
		otelmetric.WithUnit("s"))
	if err != nil {
		return janitorTelemetry{}, err
	}

	return janitorTelemetry{purgedMessages: purgedMessages, purgeRuns: purgeRuns, oldestPendingAge: oldestPendingAge}, nil
}

func (t janitorTelemetry) registerOldestPendingAge(
	read func(context.Context) (time.Duration, error),
) (otelmetric.Registration, error) {
	return meter.RegisterCallback(func(ctx context.Context, observer otelmetric.Observer) error {
		age, err := read(ctx)
		if err != nil {
			logger.ErrorBySeverity(ctx, err)

			return nil
		}

		observer.ObserveInt64(t.oldestPendingAge, int64(age.Seconds()))

		return nil
	}, t.oldestPendingAge)
}

func (janitorTelemetry) onPurgeStart(ctx context.Context) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "PurgeOutbox")

	logger.Debug(ctx, "purging outbox")

	return ctx, span
}

func (t janitorTelemetry) onPurged(ctx context.Context, purged int64) {
	t.purgedMessages.Add(ctx, purged)
	t.purgeRuns.Add(ctx, 1, outcomeSuccess)
	logger.Info(ctx, "outbox purged")
}

func (t janitorTelemetry) onPurgeFailed(ctx context.Context, purged int64, err error) error {
	t.purgedMessages.Add(ctx, purged)
	t.purgeRuns.Add(ctx, 1, outcomeFailure)

	logger.ErrorBySeverity(ctx, err)

	return err
}

func onForwarderStarted(ctx context.Context, name string, topics []string) {
	logger.Info(ctx, "forwarder started", log.String("outbox", name), log.String("topics", strings.Join(topics, ",")))
}

func onForwarderStoppedUnexpectedly(ctx context.Context, err error) {
	logger.Error(ctx, ErrForwarderNotRunning.WithCause(err))
}

func messageAttributes(msg *message.Message) []attribute.KeyValue {
	attributes := []attribute.KeyValue{semconv.MessagingMessageID(msg.UUID)}

	if key := msg.Metadata.Get(PartitionKeyHeader); key != "" {
		attributes = append(attributes, semconv.MessagingKafkaMessageKey(key))
	}

	return attributes
}

func onCreateStart(ctx context.Context, msg *message.Message, topic string) func(error) {
	ctx, span := tracer.Start(ctx, operationCreate+" "+topic,
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.MessagingSystemKafka,
			semconv.MessagingOperationTypeCreate,
			semconv.MessagingOperationName(operationCreate),
			semconv.MessagingDestinationName(topic),
		),
		trace.WithAttributes(messageAttributes(msg)...),
	)

	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Metadata))
	msg.SetContext(ctx)

	return func(err error) {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}
}

func sendAttributes(topic string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKafka,
		semconv.MessagingDestinationName(topic),
		semconv.MessagingOperationName(operationSend),
		semconv.MessagingOperationTypeSend,
	}
}

func deliveryAttempt(n int) attribute.KeyValue {
	return attribute.Int("messaging.delivery.attempt", n)
}

func onSendStart(msg *message.Message, outboxName, topic string, attempt int) func(error) {
	started := time.Now()

	attributes := append(sendAttributes(topic), deliveryAttempt(attempt))
	attributes = append(attributes, messageAttributes(msg)...)

	opts := []trace.SpanStartOption{
		trace.WithNewRoot(),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(attributes...),
	}

	producer := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(msg.Metadata))
	if link := trace.LinkFromContext(producer); link.SpanContext.IsValid() {
		opts = append(opts, trace.WithLinks(link))
	}

	ctx, span := tracer.Start(msg.Context(), operationSend+" "+topic, opts...)

	return func(err error) {
		defer span.End()

		metricAttributes := append(sendAttributes(topic), attribute.String("outbox.name", outboxName))

		if err != nil {
			errType := semconv.ErrorTypeKey.String(errorType(err))
			metricAttributes = append(metricAttributes, errType)

			span.SetAttributes(errType)
			logger.ErrorBySeverity(ctx, ErrSend.WithCause(err).WithAttribute("messaging.delivery.attempt", strconv.Itoa(attempt)))
		}

		sentMessages.Add(ctx, 1, otelmetric.WithAttributes(metricAttributes...))
		operationDuration.Record(ctx, time.Since(started).Seconds(), otelmetric.WithAttributes(metricAttributes...))
	}
}

func errorType(err error) string {
	var kafkaErr sarama.KError
	if errors.As(err, &kafkaErr) {
		return "KAFKA_ERROR_" + strconv.Itoa(int(kafkaErr))
	}

	if code := errors.Code(err); code != errors.CodeUnknown {
		return string(code)
	}

	return semconv.ErrorTypeOther.Value.AsString()
}
