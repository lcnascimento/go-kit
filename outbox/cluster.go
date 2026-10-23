package outbox

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

const (
	connectTimeout = 10 * time.Second
	pingTimeout    = time.Second
)

// cluster is the forwarder's one metadata connection: it checks the topics at start and answers readiness.
type cluster struct {
	client sarama.Client
	admin  sarama.ClusterAdmin
	topics []string
}

func dialCluster(ctx context.Context, target KafkaTarget) (*cluster, error) {
	cfg := sarama.NewConfig()
	cfg.Net.DialTimeout = pingTimeout
	cfg.Net.ReadTimeout = pingTimeout
	cfg.Net.WriteTimeout = pingTimeout
	cfg.Metadata.Retry.Max = 0
	cfg.Metadata.Timeout = pingTimeout
	cfg.Metadata.AllowAutoTopicCreation = false

	client, err := withContext(ctx, func() (sarama.Client, error) { return sarama.NewClient(target.Brokers, cfg) },
		func(c sarama.Client) { _ = c.Close() })
	if err != nil {
		return nil, ErrKafkaUnavailable.WithCause(err)
	}

	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		_ = client.Close()

		return nil, ErrKafkaUnavailable.WithCause(err)
	}

	return &cluster{client: client, admin: admin, topics: target.Topics}, nil
}

func (c *cluster) checkTopics(ctx context.Context) error {
	existing, err := withContext(ctx, c.admin.ListTopics, nil)
	if err != nil {
		return ErrKafkaUnavailable.WithCause(err)
	}

	for _, topic := range c.topics {
		if _, ok := existing[topic]; !ok {
			return ErrTopicMissing.WithAttribute("topic", topic)
		}
	}

	return nil
}

// Ping asks only for the metadata of the forwarder's topics, since a leaderless partition elsewhere
// (__consumer_offsets while the broker starts) does not stop it.
func (c *cluster) Ping(ctx context.Context) error {
	// A port that accepts but never answers stalls sarama for one read timeout per request of the handshake.
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	_, err := withContext(ctx, func() (struct{}, error) { return struct{}{}, c.client.RefreshMetadata(c.topics...) }, nil)
	if err != nil {
		return ErrKafkaUnavailable.WithCause(err)
	}

	return nil
}

// Close closes the admin, which closes the client it was made from.
func (c *cluster) Close() error {
	return c.admin.Close()
}

// withContext bounds a blocking sarama call by ctx and connectTimeout; release frees a result that arrives too late.
func withContext[T any](ctx context.Context, call func() (T, error), release func(T)) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	type result struct {
		value T
		err   error
	}

	done := make(chan result, 1)

	go func() {
		value, err := call()
		done <- result{value: value, err: err}
	}()

	select {
	case r := <-done:
		return r.value, r.err
	case <-ctx.Done():
		go func() {
			if r := <-done; r.err == nil && release != nil {
				release(r.value)
			}
		}()

		var zero T

		return zero, ctx.Err()
	}
}
