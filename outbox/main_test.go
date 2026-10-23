package outbox_test

import (
	"embed"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

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

	os.Exit(m.Run())
}

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if testing.Short() {
		t.Skip("component test: needs Postgres")
	}

	return testdb.New(t, migrations)
}

func count(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()

	var n int

	require.NoError(t, pool.QueryRow(t.Context(), `SELECT count(*) FROM watermill_test_outbox`).Scan(&n))

	return n
}

func mustSub(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}

	return sub
}
