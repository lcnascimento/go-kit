package outbox

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

const pingTimeout = time.Second

type cluster struct {
	client sarama.Client
	topics []string
}

// dialCluster connects to the brokers; Ping then asks only for the metadata of the forwarder's topics, since a
// leaderless partition elsewhere (__consumer_offsets while the broker starts) does not stop it.
func dialCluster(brokers, topics []string) (*cluster, error) {
	cfg := sarama.NewConfig()
	cfg.Net.DialTimeout = pingTimeout
	cfg.Net.ReadTimeout = pingTimeout
	cfg.Net.WriteTimeout = pingTimeout
	cfg.Metadata.Retry.Max = 0
	cfg.Metadata.Timeout = pingTimeout

	client, err := sarama.NewClient(brokers, cfg)
	if err != nil {
		return nil, ErrKafkaUnavailable.WithCause(err)
	}

	return &cluster{client: client, topics: topics}, nil
}

func (c *cluster) Ping(ctx context.Context) error {
	// A port that accepts but never answers stalls sarama for one read timeout per request of the handshake.
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	refreshed := make(chan error, 1)

	go func() { refreshed <- c.client.RefreshMetadata(c.topics...) }()

	select {
	case err := <-refreshed:
		if err != nil {
			return ErrKafkaUnavailable.WithCause(err)
		}

		return nil
	case <-ctx.Done():
		return ErrKafkaUnavailable.WithCause(ctx.Err())
	}
}

func (c *cluster) Close() error {
	return c.client.Close()
}
