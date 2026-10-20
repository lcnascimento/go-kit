package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/http/httpserver/middlewares"
)

var spans = tracetest.NewSpanRecorder()

func TestMain(m *testing.M) {
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans)))
	otel.SetTextMapPropagator(propagation.TraceContext{})

	os.Exit(m.Run())
}

func serve(t *testing.T, req *http.Request) sdktrace.ReadOnlySpan {
	t.Helper()

	before := len(spans.Ended())

	router := mux.NewRouter()
	router.Use(middlewares.Telemetry)
	router.HandleFunc("/notes/{id}", func(http.ResponseWriter, *http.Request) {})

	router.ServeHTTP(httptest.NewRecorder(), req)

	ended := spans.Ended()
	require.Len(t, ended, before+1)

	return ended[before]
}

func TestTelemetry_ContinuesTheCallerTrace(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/notes/42", nil)
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

	span := serve(t, req)

	require.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", span.SpanContext().TraceID().String())
	require.Equal(t, "00f067aa0ba902b7", span.Parent().SpanID().String())
	require.True(t, span.Parent().IsRemote())
}

func TestTelemetry_OpensAServerSpanNamedByMethodAndRoute(t *testing.T) {
	span := serve(t, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/notes/42", nil))

	require.Equal(t, trace.SpanKindServer, span.SpanKind())
	require.Equal(t, "GET /notes/{id}", span.Name())
}
