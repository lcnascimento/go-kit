package outbox

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	startTimeout = 30 * time.Second
	closeTimeout = 30 * time.Second

	retryMaxRetries      = 10
	retryInitialInterval = 100 * time.Millisecond
	retryMaxInterval     = 30 * time.Second
	retryMultiplier      = 2
)

// Forwarder sends the outbox messages to Kafka in commit order, keyed by PartitionKeyHeader, with idempotent writes.
type Forwarder struct {
	outbox    Outbox
	topics    []string
	router    *message.Router
	publisher *kafka.Publisher
	cluster   *cluster
	shutdown  func()

	started  atomic.Bool
	stopping atomic.Bool
	stopOnce sync.Once
	stopErr  error
	runDone  chan struct{}
}

// NewForwarder fails unless every topic the outbox forwards to exists.
func NewForwarder(o Outbox, pool *pgxpool.Pool, brokers, topics []string, opts ...ForwarderOption) (*Forwarder, error) {
	if pool == nil || len(brokers) == 0 || len(topics) == 0 {
		return nil, ErrInvalidForwarderConfig
	}

	if err := checkTopics(brokers, topics); err != nil {
		return nil, err
	}

	subscriber, err := o.newSubscriber(pool)
	if err != nil {
		return nil, ErrInvalidForwarderConfig.WithCause(err)
	}

	cluster, err := dialCluster(brokers, topics)
	if err != nil {
		return nil, err
	}

	publisher, err := newKafkaPublisher(brokers)
	if err != nil {
		_ = cluster.Close()

		return nil, err
	}

	f := &Forwarder{
		outbox:    o,
		topics:    topics,
		publisher: publisher,
		cluster:   cluster,
		shutdown:  func() {},
		runDone:   make(chan struct{}),
	}
	for _, opt := range opts {
		opt(f)
	}

	f.router, err = newRouter(o, subscriber, publisher)
	if err != nil {
		f.close()

		return nil, ErrInvalidForwarderConfig.WithCause(err)
	}

	return f, nil
}

func newRouter(o Outbox, subscriber message.Subscriber, publisher message.Publisher) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{CloseTimeout: closeTimeout}, watermill.NopLogger{})
	if err != nil {
		return nil, err
	}

	retry := middleware.Retry{
		MaxRetries:      retryMaxRetries,
		InitialInterval: retryInitialInterval,
		MaxInterval:     retryMaxInterval,
		Multiplier:      retryMultiplier,
	}

	_, err = forwarder.NewForwarder(subscriber, newTracedPublisher(publisher), watermill.NopLogger{}, forwarder.Config{
		ForwarderTopic: o.forwarderTopic(),
		Router:         router,
		Middlewares:    []message.HandlerMiddleware{retry.Middleware, middleware.Recoverer},
	})
	if err != nil {
		return nil, err
	}

	return router, nil
}

func (f *Forwarder) Start(ctx context.Context) error {
	runErr := make(chan error, 1)

	f.started.Store(true)

	runCtx := context.WithoutCancel(ctx)

	go func() {
		defer close(f.runDone)

		err := f.router.Run(runCtx)
		runErr <- err

		if !f.stopping.Load() {
			onForwarderStoppedUnexpectedly(runCtx, err)
			f.shutdown()
		}
	}()

	timeout := time.NewTimer(startTimeout)
	defer timeout.Stop()

	select {
	case <-f.router.Running():
		onForwarderStarted(ctx, f.outbox.name, f.topics)

		return nil
	case err := <-runErr:
		return ErrStartForwarder.WithCause(err)
	case <-timeout.C:
		return ErrStartForwarder.WithCause(context.DeadlineExceeded)
	case <-ctx.Done():
		return ErrStartForwarder.WithCause(ctx.Err())
	}
}

func (f *Forwarder) Stop(context.Context) error {
	f.stopOnce.Do(func() {
		f.stopping.Store(true)

		routerErr := f.router.Close()
		if f.started.Load() {
			<-f.runDone
		}

		publisherErr, clusterErr := f.publisher.Close(), f.cluster.Close()

		for _, err := range []error{routerErr, publisherErr, clusterErr} {
			if err != nil {
				f.stopErr = ErrStopForwarder.WithCause(err)

				return
			}
		}
	})

	return f.stopErr
}

// Ready fails while the forwarder is not running or Kafka does not answer for its topics.
func (f *Forwarder) Ready(ctx context.Context) error {
	if !f.router.IsRunning() || f.router.IsClosed() {
		return ErrForwarderNotRunning
	}

	return f.cluster.Ping(ctx)
}

func (f *Forwarder) close() {
	_ = f.publisher.Close()
	_ = f.cluster.Close()
}

func newKafkaPublisher(brokers []string) (*kafka.Publisher, error) {
	cfg := kafka.DefaultSaramaSyncPublisherConfig()
	cfg.Producer.Idempotent = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Net.MaxOpenRequests = 1

	publisher, err := kafka.NewPublisher(kafka.PublisherConfig{
		Brokers:               brokers,
		Marshaler:             kafka.NewWithPartitioningMarshaler(partitionKey),
		OverwriteSaramaConfig: cfg,
	}, watermill.NopLogger{})
	if err != nil {
		return nil, ErrKafkaUnavailable.WithCause(err)
	}

	return publisher, nil
}

func partitionKey(_ string, msg *message.Message) (string, error) {
	return msg.Metadata.Get(PartitionKeyHeader), nil
}

func checkTopics(brokers, topics []string) error {
	admin, err := sarama.NewClusterAdmin(brokers, sarama.NewConfig())
	if err != nil {
		return ErrKafkaUnavailable.WithCause(err)
	}

	defer func() { _ = admin.Close() }()

	existing, err := admin.ListTopics()
	if err != nil {
		return ErrKafkaUnavailable.WithCause(err)
	}

	for _, topic := range topics {
		if _, ok := existing[topic]; !ok {
			return ErrTopicMissing.WithAttribute("topic", topic)
		}
	}

	return nil
}
