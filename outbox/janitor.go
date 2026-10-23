package outbox

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	otelmetric "go.opentelemetry.io/otel/metric"
)

const (
	defaultRetention     = 7 * 24 * time.Hour
	defaultPurgeInterval = time.Hour

	purgeBatchSize = 10_000
)

var runningJanitors sync.Map

// Janitor purges forwarded messages older than the retention and reports the age of the oldest pending one.
type Janitor struct {
	outbox    Outbox
	pool      *pgxpool.Pool
	retention time.Duration
	interval  time.Duration
	telemetry janitorTelemetry

	lifecycle    lifecycle
	cancel       context.CancelFunc
	done         chan struct{}
	registration otelmetric.Registration
}

func NewJanitor(o Outbox, pool *pgxpool.Pool, opts ...JanitorOption) (*Janitor, error) {
	telemetry, err := newJanitorTelemetry(o.name)
	if err != nil {
		return nil, ErrNewJanitor.WithCause(err)
	}

	j := &Janitor{
		outbox:    o,
		pool:      pool,
		retention: defaultRetention,
		interval:  defaultPurgeInterval,
		telemetry: telemetry,
	}

	for _, opt := range opts {
		opt(j)
	}

	return j, nil
}

func (j *Janitor) Start(ctx context.Context) error {
	if err := j.lifecycle.start(); err != nil {
		return err
	}

	if _, taken := runningJanitors.LoadOrStore(j.outbox.name, j); taken {
		j.lifecycle.stop()

		return ErrJanitorAlreadyRunning.WithAttribute("outbox", j.outbox.name)
	}

	registration, err := j.telemetry.registerOldestPendingAge(j.OldestPendingAge)
	if err != nil {
		j.lifecycle.stop()
		runningJanitors.Delete(j.outbox.name)

		return ErrStartJanitor.WithCause(err)
	}

	j.registration = registration

	ctx, j.cancel = context.WithCancel(context.WithoutCancel(ctx))
	j.done = make(chan struct{})

	go j.run(ctx)

	return nil
}

// Stop stops reporting the age at once and waits for an ongoing purge until ctx expires.
func (j *Janitor) Stop(ctx context.Context) error {
	if j.lifecycle.stop() != running {
		return nil
	}

	defer runningJanitors.Delete(j.outbox.name)

	j.cancel()

	unregisterErr := j.registration.Unregister()

	select {
	case <-j.done:
	case <-ctx.Done():
		return ErrStopJanitor.WithCause(ctx.Err())
	}

	if unregisterErr != nil {
		return ErrStopJanitor.WithCause(unregisterErr)
	}

	return nil
}

func (j *Janitor) Purge(ctx context.Context) (int64, error) {
	ctx, span := j.telemetry.onPurgeStart(ctx)
	defer span.End()

	query, args := j.outbox.purgeQuery(j.retention, purgeBatchSize)

	var total int64

	for {
		tag, err := j.pool.Exec(ctx, query, args...)
		if err != nil {
			return total, j.telemetry.onPurgeFailed(ctx, total, ErrPurge.WithCause(err))
		}

		total += tag.RowsAffected()

		if tag.RowsAffected() < purgeBatchSize {
			j.telemetry.onPurged(ctx, total)

			return total, nil
		}
	}
}

// OldestPendingAge is 0 when every message was forwarded.
func (j *Janitor) OldestPendingAge(ctx context.Context) (time.Duration, error) {
	query, args := j.outbox.oldestPendingAgeQuery()

	var seconds float64

	if err := j.pool.QueryRow(ctx, query, args...).Scan(&seconds); err != nil {
		return 0, ErrOldestPendingAge.WithCause(err)
	}

	return time.Duration(seconds * float64(time.Second)), nil
}

func (j *Janitor) run(ctx context.Context) {
	defer close(j.done)

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		_, _ = j.Purge(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
