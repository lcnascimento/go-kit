package outbox

import (
	"context"
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

// KafkaTarget is where a Forwarder sends: the brokers, and every topic the outbox's messages go to.
type KafkaTarget struct {
	Brokers []string
	Topics  []string
}

// Forwarder sends the outbox messages to Kafka in commit order, keyed by PartitionKeyHeader, with idempotent writes.
type Forwarder struct {
	outbox    Outbox
	target    KafkaTarget
	router    *message.Router
	publisher *kafka.Publisher
	cluster   *cluster
	shutdown  func()
	validate  func(*message.Message) error

	lifecycle lifecycle
	runDone   chan struct{}
	stopDone  chan struct{}
	stopErr   error
}

// NewForwarder fails unless Kafka answers within ctx and every topic of target exists.
func NewForwarder(ctx context.Context, o Outbox, pool *pgxpool.Pool, target KafkaTarget, opts ...ForwarderOption) (*Forwarder, error) {
	if pool == nil || len(target.Brokers) == 0 || len(target.Topics) == 0 {
		return nil, ErrInvalidForwarderConfig
	}

	subscriber, err := o.newSubscriber(pool)
	if err != nil {
		return nil, ErrInvalidForwarderConfig.WithCause(err)
	}

	cluster, err := dialCluster(ctx, target)
	if err != nil {
		return nil, err
	}

	if err := cluster.checkTopics(ctx); err != nil {
		_ = cluster.Close()

		return nil, err
	}

	publisher, err := newKafkaPublisher(ctx, target.Brokers)
	if err != nil {
		_ = cluster.Close()

		return nil, err
	}

	f := &Forwarder{
		outbox:    o,
		target:    target,
		publisher: publisher,
		cluster:   cluster,
		shutdown:  func() {},
		runDone:   make(chan struct{}),
		stopDone:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(f)
	}

	f.router, err = newRouter(o, subscriber, newTracedPublisher(publisher, o.name, f.validate))
	if err != nil {
		_ = f.closeClients()

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

	_, err = forwarder.NewForwarder(subscriber, publisher, watermill.NopLogger{}, forwarder.Config{
		ForwarderTopic: o.forwarderTopic(),
		Router:         router,
		Middlewares:    []message.HandlerMiddleware{retry.Middleware, middleware.Recoverer},
	})
	if err != nil {
		return nil, err
	}

	return router, nil
}

// Start fails, releasing everything, unless the forwarder is running within ctx and startTimeout.
func (f *Forwarder) Start(ctx context.Context) error {
	if err := f.lifecycle.start(); err != nil {
		return err
	}

	runErr := make(chan error, 1)
	runCtx := context.WithoutCancel(ctx)

	go func() {
		defer close(f.runDone)

		err := f.router.Run(runCtx)
		runErr <- err

		if f.lifecycle.is(running) {
			onForwarderStoppedUnexpectedly(runCtx, err)
			f.shutdown()
		}
	}()

	timeout := time.NewTimer(startTimeout)
	defer timeout.Stop()

	select {
	case <-f.router.Running():
		onForwarderStarted(ctx, f.outbox.name, f.target.Topics)

		return nil
	case err := <-runErr:
		return f.failStart(ctx, err)
	case <-timeout.C:
		return f.failStart(ctx, context.DeadlineExceeded)
	case <-ctx.Done():
		return f.failStart(ctx, ctx.Err())
	}
}

func (f *Forwarder) failStart(ctx context.Context, cause error) error {
	err := ErrStartForwarder.WithCause(cause)

	if stopErr := f.Stop(context.WithoutCancel(ctx)); stopErr != nil {
		err = err.WithCause(stopErr)
	}

	return err
}

// Stop closes the router and the Kafka clients; it returns when ctx expires, and a later Stop waits for the rest.
func (f *Forwarder) Stop(ctx context.Context) error {
	if previous := f.lifecycle.stop(); previous != stopped {
		go f.release(previous == running)
	}

	select {
	case <-f.stopDone:
		return f.stopErr
	case <-ctx.Done():
		return ErrStopForwarder.WithCause(ctx.Err())
	}
}

func (f *Forwarder) release(wasRunning bool) {
	defer close(f.stopDone)

	var routerErr error

	if wasRunning {
		routerErr = f.router.Close()
		<-f.runDone
	}

	for _, err := range []error{routerErr, f.closeClients()} {
		if err != nil {
			f.stopErr = ErrStopForwarder.WithCause(err)

			return
		}
	}
}

// Ready fails while the forwarder is not running or Kafka does not answer for its topics.
func (f *Forwarder) Ready(ctx context.Context) error {
	if !f.lifecycle.is(running) || !f.router.IsRunning() || f.router.IsClosed() {
		return ErrForwarderNotRunning
	}

	return f.cluster.Ping(ctx)
}

func (f *Forwarder) closeClients() error {
	publisherErr, clusterErr := f.publisher.Close(), f.cluster.Close()
	if publisherErr != nil {
		return publisherErr
	}

	return clusterErr
}

// newKafkaPublisher opens the connection that sends; watermill-kafka dials its own producer and cannot share cluster's.
func newKafkaPublisher(ctx context.Context, brokers []string) (*kafka.Publisher, error) {
	cfg := kafka.DefaultSaramaSyncPublisherConfig()
	cfg.Producer.Idempotent = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Net.MaxOpenRequests = 1
	cfg.Net.DialTimeout = connectTimeout

	publisher, err := withContext(ctx, func() (*kafka.Publisher, error) {
		return kafka.NewPublisher(kafka.PublisherConfig{
			Brokers:               brokers,
			Marshaler:             kafka.NewWithPartitioningMarshaler(partitionKey),
			OverwriteSaramaConfig: cfg,
		}, watermill.NopLogger{})
	}, func(p *kafka.Publisher) { _ = p.Close() })
	if err != nil {
		return nil, ErrKafkaUnavailable.WithCause(err)
	}

	return publisher, nil
}

func partitionKey(_ string, msg *message.Message) (string, error) {
	return msg.Metadata.Get(PartitionKeyHeader), nil
}
