package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type testEvent struct {
	ID string `json:"id"`
}

func (e *testEvent) GetTopic() string         { return "" }
func (e *testEvent) GetKey() []byte           { return []byte(e.ID) }
func (e *testEvent) GetType() EventType       { return EventType("TEST") }
func (e *testEvent) GetVersion() EventVersion { return EventVersion("v1") }

// Messages published while no consumer is running must still be delivered
// once the group subscribes: subscribers promise at-least-once, and a
// LastOffset start silently drops everything already in the partition.
func TestComponentSubscriberConsumesMessagesProducedBeforeStart(t *testing.T) {
	if testing.Short() {
		t.Skip("component test requires a local Kafka broker")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	topic := fmt.Sprintf("gokit.test.start-offset.%s", uuid.NewString())
	if err := EnsureTopics(ctx, TopicConfig{Name: topic, NumPartitions: 1, ReplicationFactor: 1}); err != nil {
		t.Fatalf("EnsureTopics: %v", err)
	}

	payload, err := json.Marshal(&testEvent{ID: "evt-1"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	writer := &kafka.Writer{Addr: kafka.TCP(brokers...), Topic: topic, RequiredAcks: kafka.RequireOne}

	// Freshly created topics take a moment to elect a partition leader.
	for {
		err = writer.WriteMessages(ctx, kafka.Message{Key: []byte("evt-1"), Value: payload})
		if err == nil {
			break
		}
		if ctx.Err() != nil {
			t.Fatalf("produce: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	prevGroup := groupID
	groupID = fmt.Sprintf("gokit-test-%s", uuid.NewString())
	t.Cleanup(func() { groupID = prevGroup })

	subscriber := NewSubscriber[*testEvent](topic)
	t.Cleanup(func() { _ = subscriber.Stop(context.Background()) })

	received := make(chan string, 1)

	go func() {
		_ = subscriber.Run(ctx, func(_ context.Context, e *testEvent) error {
			received <- e.ID
			return nil
		})
	}()

	select {
	case id := <-received:
		if id != "evt-1" {
			t.Fatalf("expected evt-1, got %s", id)
		}
	case <-time.After(45 * time.Second):
		t.Fatal("message produced before the subscriber started was never consumed")
	}
}
