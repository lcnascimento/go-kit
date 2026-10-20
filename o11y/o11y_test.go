package o11y_test

import (
	"net/http"
	"net/http/httptest"
	"runtime/pprof"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"

	"github.com/lcnascimento/go-kit/o11y"
)

func TestStart_PyroscopeProfilerPushesProfilesAndLinksRootSpans(t *testing.T) {
	var (
		mu    sync.Mutex
		names []string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.URL.Path == "/ingest" {
			names = append(names, r.URL.Query().Get("name"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	t.Setenv("OTEL_SERVICE_NAME", "svc")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=ns,deployment.environment=local")
	t.Setenv("OTEL_TRACES_EXPORTER", "stdout")
	t.Setenv("OTEL_METRICS_EXPORTER", "none")
	t.Setenv("OTEL_LOGS_EXPORTER", "none")
	t.Setenv("OTEL_PROFILES_EXPORTER", "pyroscope")
	t.Setenv("OTEL_PROFILES", "cpu")
	t.Setenv("OTEL_EXPORTER_PYROSCOPE_ENDPOINT", server.URL)

	require.NotPanics(t, func() { o11y.MustStart(t.Context()) })

	ctx, span := otel.Tracer("test").Start(t.Context(), "root")
	spanID, ok := pprof.Label(ctx, "span_id")
	require.True(t, ok)
	require.Equal(t, span.SpanContext().SpanID().String(), spanID)
	span.End()

	require.NoError(t, o11y.Shutdown(t.Context()))

	mu.Lock()
	defer mu.Unlock()
	require.NotEmpty(t, names)
	require.Regexp(t, `^svc\{`, names[0])
	require.Contains(t, names[0], "service_namespace=ns")
	require.Contains(t, names[0], "deployment_environment=local")
}
