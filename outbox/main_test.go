package outbox_test

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/goleak"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/lcnascimento/go-kit/postgres/testdb"

	"github.com/lcnascimento/go-kit/outbox"
)

const waitFor = 30 * time.Second

var (
	spans   = tracetest.NewSpanRecorder()
	metrics = sdkmetric.NewManualReader()

	//go:embed testdata/migrations/*.sql
	testdata   embed.FS
	migrations = mustSub(testdata, "testdata/migrations")

	testOutbox = outbox.MustNew("test")
)

func TestMain(m *testing.M) {
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans)))
	otel.SetMeterProvider(sdkmetric.NewMeterProvider(sdkmetric.WithReader(metrics)))
	otel.SetTextMapPropagator(propagation.TraceContext{})

	goleak.VerifyTestMain(m,
		// testcontainers keeps its Docker connections and reaper for the run, and sarama starts a process-wide metrics ticker.
		goleak.IgnoreAnyFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreAnyFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("github.com/rcrowley/go-metrics.(*meterArbiter).tick"),
		goleak.IgnoreTopFunction("github.com/testcontainers/testcontainers-go.(*Reaper).connect.func1"),
	)
}

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if testing.Short() {
		t.Skip("component test: needs Postgres")
	}

	return testdb.New(t, migrations)
}

func newOutbox(t *testing.T, pool *pgxpool.Pool) outbox.Outbox {
	t.Helper()

	name := withoutTables().Name()

	raw, err := fs.ReadFile(migrations, "00001_create_test_outbox.sql")
	require.NoError(t, err)

	up, _, _ := strings.Cut(strings.TrimPrefix(string(raw), "-- +goose Up"), "-- +goose Down")

	_, err = pool.Exec(t.Context(), strings.ReplaceAll(up, "test_outbox", name+"_outbox"))
	require.NoError(t, err)

	return outbox.MustNew(name)
}

func withoutTables() outbox.Outbox {
	suffix := make([]byte, 4)
	_, _ = rand.Read(suffix)

	return outbox.MustNew("t" + hex.EncodeToString(suffix))
}

func newJanitor(t *testing.T, o outbox.Outbox, pool *pgxpool.Pool, opts ...outbox.JanitorOption) *outbox.Janitor {
	t.Helper()

	janitor, err := outbox.NewJanitor(o, pool, opts...)
	require.NoError(t, err)

	return janitor
}

func mustSub(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}

	return sub
}
