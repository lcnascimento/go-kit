package outbox

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	otelmetric "go.opentelemetry.io/otel/metric"
)

const (
	defaultRetention     = 7 * 24 * time.Hour
	defaultPurgeInterval = time.Hour

	purgeBatchSize = 10_000
)

// Janitor purges forwarded messages older than the retention and reports the age of the oldest pending one.
type Janitor struct {
	outbox    Outbox
	pool      *pgxpool.Pool
	retention time.Duration
	interval  time.Duration
	telemetry janitorTelemetry

	cancel       context.CancelFunc
	done         chan struct{}
	registration otelmetric.Registration
}

func NewJanitor(o Outbox, pool *pgxpool.Pool, opts ...JanitorOption) *Janitor {
	j := &Janitor{
		outbox:    o,
		pool:      pool,
		retention: defaultRetention,
		interval:  defaultPurgeInterval,
		telemetry: newJanitorTelemetry(o.name),
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

func (j *Janitor) Start(ctx context.Context) error {
	registration, err := j.telemetry.registerOldestPendingAge(j.OldestPendingAge)
	if err != nil {
		return ErrStartJanitor.WithCause(err)
	}

	j.registration = registration

	ctx, j.cancel = context.WithCancel(context.WithoutCancel(ctx))
	j.done = make(chan struct{})

	go j.run(ctx)

	return nil
}

func (j *Janitor) Stop(ctx context.Context) error {
	j.cancel()

	select {
	case <-j.done:
	case <-ctx.Done():
		return ctx.Err()
	}

	return j.registration.Unregister()
}

func (j *Janitor) Purge(ctx context.Context) (int64, error) {
	ctx, span := j.telemetry.onPurgeStart(ctx)
	defer span.End()

	var total int64

	for {
		tag, err := j.pool.Exec(ctx, `
			DELETE FROM `+j.outbox.messagesTable()+`
			WHERE (transaction_id, "offset") IN (
				SELECT m.transaction_id, m."offset"
				FROM `+j.outbox.messagesTable()+` m
				JOIN `+j.outbox.offsetsTable()+` o ON o.consumer_group = $1
				WHERE m.created_at < now() - make_interval(secs => $2::float8)
				  AND (m.transaction_id, m."offset") <= (o.last_processed_transaction_id, o.offset_acked)
				LIMIT $3
			)`,
			j.outbox.consumerGroup(), j.retention.Seconds(), purgeBatchSize)
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
	var seconds float64

	err := j.pool.QueryRow(ctx, `
		SELECT coalesce(extract(epoch FROM now() - min(m.created_at)), 0)::float8
		FROM `+j.outbox.messagesTable()+` m
		LEFT JOIN `+j.outbox.offsetsTable()+` o ON o.consumer_group = $1
		WHERE o.offset_acked IS NULL
		   OR (m.transaction_id, m."offset") > (o.last_processed_transaction_id, o.offset_acked)`,
		j.outbox.consumerGroup()).Scan(&seconds)
	if err != nil {
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
