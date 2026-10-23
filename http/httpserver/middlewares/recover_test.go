package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/lcnascimento/go-kit/o11y"

	"github.com/lcnascimento/go-kit/http/httpserver/middlewares"
)

func TestRecover_CapturesThePanicAsAnException(t *testing.T) {
	var (
		mu      sync.Mutex
		batches []map[string]any
	)
	posthog := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Batch []map[string]any `json:"batch"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
			mu.Lock()
			batches = append(batches, body.Batch...)
			mu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(posthog.Close)

	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	t.Setenv("OTEL_METRICS_EXPORTER", "none")
	t.Setenv("OTEL_LOGS_EXPORTER", "none")
	t.Setenv("POSTHOG_PROJECT_KEY", "phc_test")
	t.Setenv("POSTHOG_HOST", posthog.URL)
	o11y.MustStart(t.Context())

	router := mux.NewRouter()
	router.Use(middlewares.Recover)
	router.HandleFunc("/explode", func(http.ResponseWriter, *http.Request) { explode() })

	res := httptest.NewRecorder()
	router.ServeHTTP(res, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/explode", nil))

	require.NoError(t, o11y.Shutdown(t.Context()))
	require.Equal(t, http.StatusInternalServerError, res.Code)

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, batches, 1)

	properties := batches[0]["properties"].(map[string]any)
	require.Equal(t, "CRITICAL", properties["kind"])

	item := properties["$exception_list"].([]any)[0].(map[string]any)
	require.Equal(t, "PANIC", item["type"])
	require.Equal(t, "kaboom", item["value"])

	inApp := map[string]bool{}
	for _, frame := range item["stacktrace"].(map[string]any)["frames"].([]any) {
		frame := frame.(map[string]any)
		inApp[frame["function"].(string)] = frame["in_app"].(bool)
	}
	require.Contains(t, inApp, "github.com/lcnascimento/go-kit/http/httpserver/middlewares_test.explode")
	require.True(t, inApp["github.com/lcnascimento/go-kit/http/httpserver/middlewares_test.explode"])
}

func explode() {
	panic("kaboom")
}
