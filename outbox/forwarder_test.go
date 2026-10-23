package outbox_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/redpanda"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/lcnascimento/go-kit/outbox"
)

const redpandaImage = "redpandadata/redpanda:v26.2.5"

var startBroker = sync.OnceValues(func() (string, error) {
	ctx := context.Background()

	container, err := redpanda.Run(ctx, redpandaImage)
	if err != nil {
		return "", err
	}

	return container.KafkaSeedBroker(ctx)
})

func TestComponentForwarder_DeliversTheMessageKeyedByItsPartitionKey(t *testing.T) {
	t.Parallel()

	env := newForwarderEnv(t)
	env.publish(t, newMessage("1f0e8a52-6f0e-7c3a-9d4b-2a1e0c9b8d7f", "note_1", "Lista de compras"))
	env.startForwarder(t)

	records := env.consume(t, "note_1", 1)

	require.Equal(t, "Lista de compras", string(records[0].Value))
	require.Equal(t, "note_1", string(records[0].Key))
	require.Equal(t, "boilerplate.notes.v1.NoteCreated", header(records[0], "ce_type"))
	require.Equal(t, "1f0e8a52-6f0e-7c3a-9d4b-2a1e0c9b8d7f", header(records[0], "_watermill_message_uuid"))
}

func TestComponentForwarder_KeepsTheOrderOfAKey(t *testing.T) {
	t.Parallel()

	env := newForwarderEnv(t)

	for _, payload := range []string{"first", "second", "third", "fourth", "fifth", "sixth"} {
		env.publish(t, newMessage(uuid.NewString(), "note_1", payload))
	}

	env.startForwarder(t)

	got := make([]string, 0, 6)
	for _, record := range env.consume(t, "note_1", 6) {
		got = append(got, string(record.Value))
	}

	require.Equal(t, []string{"first", "second", "third", "fourth", "fifth", "sixth"}, got)
}

func TestComponentForwarder_RestartLeavesNoGap(t *testing.T) {
	t.Parallel()

	env := newForwarderEnv(t)

	for range 20 {
		env.publish(t, newMessage(uuid.NewString(), "restart", "before"))
	}

	first := env.startForwarder(t)
	env.consume(t, "restart", 1)
	require.NoError(t, first.Stop(t.Context()))

	for range 20 {
		env.publish(t, newMessage(uuid.NewString(), "restart", "after"))
	}

	env.startForwarder(t)

	require.Len(t, env.consume(t, "restart", 40), 40)
}

func TestComponentForwarder_SendSpanIsARootLinkedToTheCreateSpan(t *testing.T) {
	t.Parallel()

	env := newForwarderEnv(t)
	messageID := uuid.NewString()

	env.publish(t, newMessage(messageID, "traced", "traced"))
	env.startForwarder(t)
	env.consume(t, "traced", 1)

	create := endedSpan(t, "create "+env.topic, messageID)
	send := endedSpan(t, "send "+env.topic, messageID)

	require.Equal(t, trace.SpanKindProducer, create.SpanKind())
	require.Equal(t, map[string]string{
		"messaging.system":            "kafka",
		"messaging.destination.name":  env.topic,
		"messaging.operation.name":    "create",
		"messaging.operation.type":    "create",
		"messaging.message.id":        messageID,
		"messaging.kafka.message.key": "traced",
	}, attributesOf(create))

	require.False(t, send.Parent().IsValid(), "send must be a root span")
	require.NotEqual(t, create.SpanContext().TraceID(), send.SpanContext().TraceID())
	require.Equal(t, trace.SpanKindProducer, send.SpanKind())
	require.Len(t, send.Links(), 1)
	require.Equal(t, create.SpanContext().SpanID(), send.Links()[0].SpanContext.SpanID())
	require.Equal(t, create.SpanContext().TraceID(), send.Links()[0].SpanContext.TraceID())
	require.Equal(t, map[string]string{
		"messaging.system":            "kafka",
		"messaging.destination.name":  env.topic,
		"messaging.operation.name":    "send",
		"messaging.operation.type":    "send",
		"messaging.message.id":        messageID,
		"messaging.kafka.message.key": "traced",
		"messaging.delivery.attempt":  "1",
	}, attributesOf(send))
}

func TestComponentForwarder_RecordsSendsAndTheirDurationInSecondsWithoutTheDeliveryAttempt(t *testing.T) {
	t.Parallel()

	env := newForwarderEnv(t)
	env.publish(t, newMessage(uuid.NewString(), "timed", "timed"))
	env.startForwarder(t)
	env.consume(t, "timed", 1)

	wantAttributes := map[string]string{
		"messaging.system":           "kafka",
		"messaging.destination.name": env.topic,
		"messaging.operation.name":   "send",
		"messaging.operation.type":   "send",
	}

	sent := dataPoint(t, "messaging.client.sent.messages", env.topic)
	require.Equal(t, wantAttributes, sent.attributes)
	require.Equal(t, int64(1), sent.count)

	duration := dataPoint(t, "messaging.client.operation.duration", env.topic)
	require.Equal(t, wantAttributes, duration.attributes)
	require.Equal(t, "s", duration.unit)
	require.Equal(t, []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10}, duration.bounds)
}

func TestComponentForwarder_JanitorPurgesWhatTheForwarderSent(t *testing.T) {
	t.Parallel()

	env := newForwarderEnv(t)
	janitor := outbox.NewJanitor(testOutbox, env.pool, outbox.WithRetention(0))

	env.publish(t, newMessage(uuid.NewString(), "purged", "sent"))

	forwarder := env.startForwarder(t)
	env.consume(t, "purged", 1)
	require.NoError(t, forwarder.Stop(t.Context()))

	env.publish(t, newMessage(uuid.NewString(), "purged", "pending"))

	purged, err := janitor.Purge(t.Context())
	require.NoError(t, err)
	require.Equal(t, int64(1), purged)
	require.Equal(t, 1, count(t, env.pool))

	age, err := janitor.OldestPendingAge(t.Context())
	require.NoError(t, err)
	require.Positive(t, age)
}

func TestComponentForwarder_NewFailsWhenATopicIsMissing(t *testing.T) {
	t.Parallel()

	env := newForwarderEnv(t)

	_, err := outbox.NewForwarder(testOutbox, env.pool, env.brokers, []string{env.topic, "missing." + uuid.NewString()})

	require.ErrorIs(t, err, outbox.ErrTopicMissing)
}

func TestComponentForwarder_ReadyOnlyWhileRunning(t *testing.T) {
	t.Parallel()

	env := newForwarderEnv(t)
	forwarder := env.startForwarder(t)

	require.NoError(t, forwarder.Ready(t.Context()))
	require.NoError(t, forwarder.Stop(t.Context()))
	require.ErrorIs(t, forwarder.Ready(t.Context()), outbox.ErrForwarderNotRunning)
}

type forwarderEnv struct {
	brokers []string
	topic   string
	pool    *pgxpool.Pool
}

func newForwarderEnv(t *testing.T) *forwarderEnv {
	t.Helper()

	pool := newPool(t)

	broker, err := startBroker()
	require.NoError(t, err, "start Redpanda")

	env := &forwarderEnv{brokers: []string{broker}, topic: "test." + uuid.NewString(), pool: pool}

	admin, err := sarama.NewClusterAdmin(env.brokers, sarama.NewConfig())
	require.NoError(t, err)

	defer func() { _ = admin.Close() }()

	require.NoError(t, admin.CreateTopic(env.topic, &sarama.TopicDetail{NumPartitions: 3, ReplicationFactor: 1}, false))

	return env
}

func (e *forwarderEnv) publish(t *testing.T, msg *message.Message) {
	t.Helper()

	tx, err := e.pool.Begin(t.Context())
	require.NoError(t, err)

	require.NoError(t, testOutbox.Publish(t.Context(), tx, e.topic, msg))
	require.NoError(t, tx.Commit(t.Context()))
}

func (e *forwarderEnv) startForwarder(t *testing.T) *outbox.Forwarder {
	t.Helper()

	f, err := outbox.NewForwarder(testOutbox, e.pool, e.brokers, []string{e.topic}, outbox.WithShutdown(func() {
		t.Errorf("forwarder stopped unexpectedly")
	}))
	require.NoError(t, err)

	require.NoError(t, f.Start(t.Context()))
	t.Cleanup(func() { _ = f.Stop(context.Background()) })

	return f
}

func (e *forwarderEnv) consume(t *testing.T, key string, n int) []*sarama.ConsumerMessage {
	t.Helper()

	consumer, err := sarama.NewConsumer(e.brokers, sarama.NewConfig())
	require.NoError(t, err)

	t.Cleanup(func() { _ = consumer.Close() })

	partitions, err := consumer.Partitions(e.topic)
	require.NoError(t, err)

	records := make(chan *sarama.ConsumerMessage, 1000)

	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(e.topic, partition, sarama.OffsetOldest)
		require.NoError(t, err)

		t.Cleanup(func() { _ = pc.Close() })

		go func() {
			for record := range pc.Messages() {
				if string(record.Key) == key {
					records <- record
				}
			}
		}()
	}

	got := make([]*sarama.ConsumerMessage, 0, n)

	for len(got) < n {
		select {
		case record := <-records:
			got = append(got, record)
		case <-time.After(waitFor):
			t.Fatalf("got %d of %d records keyed %s", len(got), n, key)
		}
	}

	return got
}

func newMessage(id, key, payload string) *message.Message {
	msg := message.NewMessage(id, []byte(payload))
	msg.Metadata.Set("ce_type", "boilerplate.notes.v1.NoteCreated")
	msg.Metadata.Set(outbox.PartitionKeyHeader, key)

	return msg
}

func header(record *sarama.ConsumerMessage, key string) string {
	for _, h := range record.Headers {
		if string(h.Key) == key {
			return string(h.Value)
		}
	}

	return ""
}

func endedSpan(t *testing.T, name, messageID string) sdktrace.ReadOnlySpan {
	t.Helper()

	var found sdktrace.ReadOnlySpan

	require.Eventually(t, func() bool {
		for _, span := range spans.Ended() {
			if span.Name() == name && attributesOf(span)["messaging.message.id"] == messageID {
				found = span

				return true
			}
		}

		return false
	}, waitFor, 50*time.Millisecond, "no ended %q span for message %s", name, messageID)

	return found
}

func attributesOf(span sdktrace.ReadOnlySpan) map[string]string {
	attributes := map[string]string{}
	for _, kv := range span.Attributes() {
		attributes[string(kv.Key)] = kv.Value.String()
	}

	return attributes
}

type point struct {
	attributes map[string]string
	unit       string
	count      int64
	bounds     []float64
}

// dataPoint finds the one data point of the metric recorded for topic, so tests on other topics do not interfere.
func dataPoint(t *testing.T, name, topic string) point {
	t.Helper()

	var found []point

	require.Eventually(t, func() bool {
		found = nil

		var collected metricdata.ResourceMetrics
		require.NoError(t, metrics.Collect(t.Context(), &collected))

		for _, scope := range collected.ScopeMetrics {
			for _, m := range scope.Metrics {
				if m.Name != name {
					continue
				}

				switch data := m.Data.(type) {
				case metricdata.Sum[int64]:
					for _, p := range data.DataPoints {
						if destination(p.Attributes) == topic {
							found = append(found, point{attributes: mapOf(p.Attributes), unit: m.Unit, count: p.Value})
						}
					}
				case metricdata.Histogram[float64]:
					for _, p := range data.DataPoints {
						if destination(p.Attributes) == topic {
							found = append(found, point{
								attributes: mapOf(p.Attributes), unit: m.Unit, bounds: p.Bounds,
							})
						}
					}
				}
			}
		}

		return len(found) > 0
	}, waitFor, 50*time.Millisecond, "no %s data point for %s", name, topic)

	require.Len(t, found, 1, "%s for %s", name, topic)

	return found[0]
}

func destination(set attribute.Set) string {
	value, _ := set.Value("messaging.destination.name")

	return value.AsString()
}

func mapOf(set attribute.Set) map[string]string {
	out := map[string]string{}
	for _, kv := range set.ToSlice() {
		out[string(kv.Key)] = kv.Value.String()
	}

	return out
}
