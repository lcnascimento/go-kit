package outbox_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/lcnascimento/go-kit/outbox"
)

const (
	expired = 8 * 24 * time.Hour
	recent  = 24 * time.Hour
)

func TestOutbox_NewRejectsANameThatIsNotAShortIdentifier(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"", "Notes", "1notes", "notes-outbox", `notes"; DROP TABLE x; --`, strings.Repeat("a", 39)} {
		_, err := outbox.New(name)

		require.ErrorIs(t, err, outbox.ErrInvalidName, name)
	}
}

func TestOutbox_NewAcceptsTheLongestNameWhoseTablesFitAPostgresIdentifier(t *testing.T) {
	t.Parallel()

	o, err := outbox.New(strings.Repeat("a", 38))

	require.NoError(t, err)
	require.Equal(t, strings.Repeat("a", 38), o.Name())
}

func TestComponentJanitor_PurgeKeepsPendingMessages(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 2, expired)
	ack(t, pool, o)
	seed(t, pool, o, 3, expired)

	purged, err := newJanitor(t, o, pool).Purge(t.Context())

	require.NoError(t, err)
	require.Equal(t, int64(2), purged)
	require.Equal(t, 3, count(t, pool, o))
}

func TestComponentJanitor_PurgeKeepsEverythingBeforeTheFirstAck(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 3, expired)

	purged, err := newJanitor(t, o, pool).Purge(t.Context())

	require.NoError(t, err)
	require.Zero(t, purged)
	require.Equal(t, 3, count(t, pool, o))
}

func TestComponentJanitor_PurgeKeepsSevenDaysByDefault(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 1, 7*24*time.Hour+time.Minute)
	seed(t, pool, o, 3, 7*24*time.Hour-time.Minute)
	ack(t, pool, o)

	purged, err := newJanitor(t, o, pool).Purge(t.Context())

	require.NoError(t, err)
	require.Equal(t, int64(1), purged)
	require.Equal(t, 3, count(t, pool, o))
}

func TestComponentJanitor_PurgeHonorsTheRetention(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 2, 3*time.Hour)
	seed(t, pool, o, 1, time.Hour)
	ack(t, pool, o)

	purged, err := newJanitor(t, o, pool, outbox.WithRetention(2*time.Hour)).Purge(t.Context())

	require.NoError(t, err)
	require.Equal(t, int64(2), purged)
	require.Equal(t, 1, count(t, pool, o))
}

func TestComponentJanitor_PurgeKeepsRecentlyForwardedMessages(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 3, recent)
	ack(t, pool, o)

	purged, err := newJanitor(t, o, pool).Purge(t.Context())

	require.NoError(t, err)
	require.Zero(t, purged)
	require.Equal(t, 3, count(t, pool, o))
}

func TestComponentJanitor_PurgeDeletesExpiredForwardedMessagesAcrossBatches(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 10_001, expired)
	seed(t, pool, o, 2, recent)
	ack(t, pool, o)

	purged, err := newJanitor(t, o, pool).Purge(t.Context())

	require.NoError(t, err)
	require.Equal(t, int64(10_001), purged)
	require.Equal(t, 2, count(t, pool, o))
}

func TestComponentJanitor_PurgeFailsWithoutTheOutboxTablesAndMarksItsSpan(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	missing := withoutTables()

	ctx, parent := otel.Tracer("test").Start(t.Context(), "test")
	_, err := newJanitor(t, missing, pool).Purge(ctx)
	parent.End()

	require.ErrorIs(t, err, outbox.ErrPurge)

	span := childSpan(t, "PurgeOutbox", parent.SpanContext().SpanID())
	require.Equal(t, codes.Error, span.Status().Code)
}

func TestComponentJanitor_OldestPendingAgeFailsWithoutTheOutboxTables(t *testing.T) {
	t.Parallel()

	_, err := newJanitor(t, withoutTables(), newPool(t)).OldestPendingAge(t.Context())

	require.ErrorIs(t, err, outbox.ErrOldestPendingAge)
}

func TestComponentJanitor_OldestPendingAgeIsZeroWithoutMessages(t *testing.T) {
	t.Parallel()

	pool := newPool(t)

	age, err := newJanitor(t, newOutbox(t, pool), pool).OldestPendingAge(t.Context())

	require.NoError(t, err)
	require.Zero(t, age)
}

func TestComponentJanitor_OldestPendingAgeIsZeroWhenEverythingWasForwarded(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 3, time.Hour)
	ack(t, pool, o)

	age, err := newJanitor(t, o, pool).OldestPendingAge(t.Context())

	require.NoError(t, err)
	require.Zero(t, age)
}

func TestComponentJanitor_OldestPendingAgeIgnoresForwardedMessages(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 1, time.Hour)
	ack(t, pool, o)
	seed(t, pool, o, 1, 10*time.Minute)
	seed(t, pool, o, 1, time.Minute)

	age, err := newJanitor(t, o, pool).OldestPendingAge(t.Context())

	require.NoError(t, err)
	require.InDelta(t, (10 * time.Minute).Seconds(), age.Seconds(), 5)
}

func TestComponentJanitor_OldestPendingAgeCountsEverythingBeforeTheFirstAck(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 1, time.Hour)

	age, err := newJanitor(t, o, pool).OldestPendingAge(t.Context())

	require.NoError(t, err)
	require.InDelta(t, time.Hour.Seconds(), age.Seconds(), 5)
}

func TestComponentJanitor_PurgesOnStart(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 3, expired)
	ack(t, pool, o)

	startJanitor(t, newJanitor(t, o, pool))

	require.Eventually(t, func() bool { return count(t, pool, o) == 0 }, waitFor, 10*time.Millisecond)
}

func TestComponentJanitor_KeepsPurgingAfterAFailure(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 1, expired)
	ack(t, pool, o)

	renamed := pgx.Identifier{"watermill_" + o.Name() + "_outbox"}.Sanitize()
	aside := pgx.Identifier{"aside_" + o.Name()}.Sanitize()

	_, err := pool.Exec(t.Context(), `ALTER TABLE `+renamed+` RENAME TO `+aside)
	require.NoError(t, err)

	startJanitor(t, newJanitor(t, o, pool, outbox.WithPurgeInterval(time.Millisecond)))

	require.Eventually(t, func() bool {
		return counter(t, o.Name()+".outbox.purge.runs", "failure") > 0
	}, waitFor, 10*time.Millisecond)

	_, err = pool.Exec(t.Context(), `ALTER TABLE `+aside+` RENAME TO `+renamed)
	require.NoError(t, err)

	require.Eventually(t, func() bool { return count(t, pool, o) == 0 }, waitFor, 10*time.Millisecond)
}

func TestComponentJanitor_StopsPurging(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	janitor := newJanitor(t, o, pool, outbox.WithPurgeInterval(time.Millisecond))

	require.NoError(t, janitor.Start(t.Context()))
	require.NoError(t, janitor.Stop(t.Context()))

	seed(t, pool, o, 3, expired)
	ack(t, pool, o)

	time.Sleep(20 * time.Millisecond)

	require.Equal(t, 3, count(t, pool, o))
}

func TestComponentJanitor_StopBeforeStartDoesNothing(t *testing.T) {
	t.Parallel()

	pool := newPool(t)

	require.NoError(t, newJanitor(t, newOutbox(t, pool), pool).Stop(t.Context()))
}

func TestComponentJanitor_StartsOnlyOnce(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	janitor := newJanitor(t, newOutbox(t, pool), pool)
	startJanitor(t, janitor)

	require.ErrorIs(t, janitor.Start(t.Context()), outbox.ErrAlreadyStarted)
}

func TestComponentJanitor_OnlyOneJanitorRunsAnOutboxPerProcess(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	first := newJanitor(t, o, pool)
	require.NoError(t, first.Start(t.Context()))

	require.ErrorIs(t, newJanitor(t, o, pool).Start(t.Context()), outbox.ErrJanitorAlreadyRunning)

	require.NoError(t, first.Stop(t.Context()))
	startJanitor(t, newJanitor(t, o, pool))
}

func TestComponentJanitor_ReportsTheOldestPendingAgeUntilStopped(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 1, 90*time.Second)

	janitor := newJanitor(t, o, pool)
	require.NoError(t, janitor.Start(t.Context()))

	got, ok := gauge(t, o.Name()+".outbox.oldest_pending.age")
	require.True(t, ok)
	require.InDelta(t, 90, got, 5)

	require.NoError(t, janitor.Stop(t.Context()))

	_, ok = gauge(t, o.Name()+".outbox.oldest_pending.age")
	require.False(t, ok)
}

func TestComponentJanitor_StopWithAnExpiredContextStillStopsReportingTheAge(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	janitor := newJanitor(t, o, pool)
	require.NoError(t, janitor.Start(t.Context()))

	expiredCtx, cancel := context.WithCancel(t.Context())
	cancel()

	if err := janitor.Stop(expiredCtx); err != nil {
		require.ErrorIs(t, err, outbox.ErrStopJanitor)
		require.ErrorIs(t, err, context.Canceled)
	}

	_, ok := gauge(t, o.Name()+".outbox.oldest_pending.age")
	require.False(t, ok)
}

func TestComponentJanitor_SkipsTheOldestPendingAgeWhenItCannotReadIt(t *testing.T) {
	t.Parallel()

	pool := newPool(t)

	missing := withoutTables()
	startJanitor(t, newJanitor(t, missing, pool))

	_, ok := gauge(t, missing.Name()+".outbox.oldest_pending.age")
	require.False(t, ok)
}

func TestComponentJanitor_PurgeCountsRunsByOutcomeAndPurgedMessages(t *testing.T) {
	t.Parallel()

	pool := newPool(t)
	o := newOutbox(t, pool)
	seed(t, pool, o, 4, expired)
	ack(t, pool, o)

	_, err := newJanitor(t, o, pool).Purge(t.Context())
	require.NoError(t, err)

	missing := withoutTables()
	_, err = newJanitor(t, missing, pool).Purge(t.Context())
	require.Error(t, err)

	require.Equal(t, int64(1), counter(t, o.Name()+".outbox.purge.runs", "success"))
	require.Equal(t, int64(4), counter(t, o.Name()+".outbox.purged.messages", ""))
	require.Equal(t, int64(1), counter(t, missing.Name()+".outbox.purge.runs", "failure"))
}

func startJanitor(t *testing.T, janitor *outbox.Janitor) {
	t.Helper()

	require.NoError(t, janitor.Start(t.Context()))
	t.Cleanup(func() { require.NoError(t, janitor.Stop(context.Background())) })
}

func messagesTable(o outbox.Outbox) string {
	return pgx.Identifier{"watermill_" + o.Name() + "_outbox"}.Sanitize()
}

func seed(t *testing.T, pool *pgxpool.Pool, o outbox.Outbox, n int, age time.Duration) {
	t.Helper()

	_, err := pool.Exec(t.Context(), `
		INSERT INTO `+messagesTable(o)+` (uuid, created_at, payload, metadata, transaction_id)
		SELECT gen_random_uuid()::text, now() - make_interval(secs => $2), '{}', '{}', pg_current_xact_id()
		FROM generate_series(1, $1)`, n, age.Seconds())
	require.NoError(t, err)
}

func ack(t *testing.T, pool *pgxpool.Pool, o outbox.Outbox) {
	t.Helper()

	_, err := pool.Exec(t.Context(), `
		INSERT INTO `+pgx.Identifier{"watermill_offsets_" + o.Name() + "_outbox"}.Sanitize()+`
			(consumer_group, offset_acked, last_processed_transaction_id)
		SELECT $1, "offset", transaction_id
		FROM `+messagesTable(o)+`
		ORDER BY transaction_id DESC, "offset" DESC
		LIMIT 1
		ON CONFLICT (consumer_group) DO UPDATE
		SET offset_acked = excluded.offset_acked, last_processed_transaction_id = excluded.last_processed_transaction_id`,
		o.Name()+"-forwarder")
	require.NoError(t, err)
}

func count(t *testing.T, pool *pgxpool.Pool, o outbox.Outbox) int {
	t.Helper()

	var n int

	require.NoError(t, pool.QueryRow(t.Context(), `SELECT count(*) FROM `+messagesTable(o)).Scan(&n))

	return n
}

func gauge(t *testing.T, name string) (int64, bool) {
	t.Helper()

	var collected metricdata.ResourceMetrics
	require.NoError(t, metrics.Collect(t.Context(), &collected))

	for _, scope := range collected.ScopeMetrics {
		for _, m := range scope.Metrics {
			if m.Name != name {
				continue
			}

			data, ok := m.Data.(metricdata.Gauge[int64])
			require.True(t, ok, "%s is a %T", name, m.Data)

			if len(data.DataPoints) == 0 {
				return 0, false
			}

			require.Len(t, data.DataPoints, 1, name)

			return data.DataPoints[0].Value, true
		}
	}

	return 0, false
}

func counter(t *testing.T, name, outcome string) int64 {
	t.Helper()

	var collected metricdata.ResourceMetrics
	require.NoError(t, metrics.Collect(t.Context(), &collected))

	for _, scope := range collected.ScopeMetrics {
		for _, m := range scope.Metrics {
			sum, ok := m.Data.(metricdata.Sum[int64])
			if m.Name != name || !ok {
				continue
			}

			for _, point := range sum.DataPoints {
				if got, _ := point.Attributes.Value(attribute.Key("outcome")); got.AsString() == outcome {
					return point.Value
				}
			}
		}
	}

	return 0
}

func childSpan(t *testing.T, name string, parent trace.SpanID) sdktrace.ReadOnlySpan {
	t.Helper()

	for _, span := range spans.Ended() {
		if span.Name() == name && span.Parent().SpanID() == parent {
			return span
		}
	}

	t.Fatalf("no ended %q span under %s", name, parent)

	return nil
}
