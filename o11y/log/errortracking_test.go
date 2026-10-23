package log_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y"
	"github.com/lcnascimento/go-kit/o11y/baggage"

	"github.com/lcnascimento/go-kit/o11y/log"
)

type posthog struct {
	mu       sync.Mutex
	requests int
	events   []map[string]any
}

func startPostHog(t *testing.T) *posthog {
	t.Helper()

	p := &posthog{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.mu.Lock()
		defer p.mu.Unlock()

		p.requests++
		if r.URL.Path == "/batch/" {
			var body struct {
				Batch []map[string]any `json:"batch"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
				p.events = append(p.events, body.Batch...)
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	t.Setenv("OTEL_SERVICE_NAME", "svc")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "deployment.environment=local")
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	t.Setenv("OTEL_METRICS_EXPORTER", "none")
	t.Setenv("OTEL_LOGS_EXPORTER", "none")
	t.Setenv("POSTHOG_PROJECT_KEY", "phc_test")
	t.Setenv("POSTHOG_HOST", server.URL)

	return p
}

func (p *posthog) exceptions(t *testing.T) []map[string]any {
	t.Helper()

	p.mu.Lock()
	defer p.mu.Unlock()

	exceptions := []map[string]any{}
	for _, event := range p.events {
		if event["event"] == "$exception" {
			exceptions = append(exceptions, event["properties"].(map[string]any))
		}
	}

	return exceptions
}

func newLogger() log.Logger {
	return log.MustNewLogger("test", log.WithLogHandler(slog.NewTextHandler(io.Discard, nil)))
}

func TestErrorTracking_CapturesAnInternalErrorAsAnException(t *testing.T) {
	posthog := startPostHog(t)
	o11y.MustStart(t.Context())

	ctx, span := trace.NewTracerProvider().Tracer("test").Start(t.Context(), "op")
	ctx = baggage.ContextWithMembers(ctx,
		baggage.NewMember("user.id", "user_1"),
		baggage.NewMember("organization.id", "org_1"),
		baggage.NewMember("tenant", "from-baggage"),
		baggage.NewMember("kind", "from-baggage"),
		baggage.NewMember("attempt", "from-baggage"),
	)

	err := errors.New("failed to save note").
		WithKind(errors.KindInternal).
		WithCode("SAVE_NOTE").
		WithAttribute("attempt", "3").
		WithAttribute("kind", "from-error").
		WithCause(errors.New("connection refused"))

	newLogger().Error(ctx, err, log.String("ignored", "attr"))
	span.End()

	require.NoError(t, o11y.Shutdown(t.Context()))

	exceptions := posthog.exceptions(t)
	require.Len(t, exceptions, 1)

	exception := exceptions[0]
	require.Equal(t, "user_1", exception["distinct_id"])
	require.Equal(t, map[string]any{"organization": "org_1"}, exception["$groups"])
	require.NotContains(t, exception, "$process_person_profile")
	require.NotContains(t, exception, "$session_id")
	require.Equal(t, span.SpanContext().TraceID().String(), exception["$trace_id"])
	require.Equal(t, span.SpanContext().SpanID().String(), exception["$span_id"])

	require.Equal(t, "INTERNAL", exception["kind"])
	require.Equal(t, "SAVE_NOTE", exception["code"])
	require.Equal(t, "ERROR", exception["severity"])
	require.Equal(t, false, exception["retryable"])
	require.Equal(t, "svc", exception["service.name"])
	require.Equal(t, "local", exception["deployment.environment"])
	require.Equal(t, []any{"connection refused"}, exception["reasons"])
	require.Equal(t, "3", exception["attempt"])
	require.Equal(t, "from-baggage", exception["tenant"])
	require.Equal(t, "user_1", exception["user.id"])
	require.Equal(t, "org_1", exception["organization.id"])
	require.NotContains(t, exception, "ignored")

	list := exception["$exception_list"].([]any)
	require.Len(t, list, 1)

	item := list[0].(map[string]any)
	require.Equal(t, "SAVE_NOTE", item["type"])
	require.Equal(t, "failed to save note", item["value"])

	frames := item["stacktrace"].(map[string]any)["frames"].([]any)
	top := frames[len(frames)-1].(map[string]any)
	require.Equal(t, "github.com/lcnascimento/go-kit/o11y/log_test.TestErrorTracking_CapturesAnInternalErrorAsAnException", top["function"])
	require.Equal(t, true, top["in_app"])
}

func TestErrorTracking_SendsTheSessionIDFromTheBaggage(t *testing.T) {
	posthog := startPostHog(t)
	o11y.MustStart(t.Context())

	ctx := baggage.ContextWithMembers(t.Context(),
		baggage.NewMember("user.id", "user_1"),
		baggage.NewMember("posthog.session_id", "0192-session"),
	)
	newLogger().Error(ctx, errors.New("boom").WithKind(errors.KindInternal))

	require.NoError(t, o11y.Shutdown(t.Context()))

	exceptions := posthog.exceptions(t)
	require.Len(t, exceptions, 1)
	require.Equal(t, "0192-session", exceptions[0]["$session_id"])
	require.NotContains(t, exceptions[0], "posthog.session_id")
}

func TestErrorTracking_IdentifiesAnAnonymousErrorByTheServiceName(t *testing.T) {
	posthog := startPostHog(t)
	o11y.MustStart(t.Context())

	newLogger().ErrorMessage(t.Context(), "boom")

	require.NoError(t, o11y.Shutdown(t.Context()))

	exceptions := posthog.exceptions(t)
	require.Len(t, exceptions, 1)
	require.Equal(t, "svc", exceptions[0]["distinct_id"])
	require.Equal(t, false, exceptions[0]["$process_person_profile"])
	require.NotContains(t, exceptions[0], "$groups")
	require.NotContains(t, exceptions[0], "$trace_id")

	item := exceptions[0]["$exception_list"].([]any)[0].(map[string]any)
	require.Equal(t, "UNKNOWN", item["type"])
	require.Equal(t, "boom", item["value"])
}

func TestErrorTracking_CapturesOnlyUnexpectedKinds(t *testing.T) {
	posthog := startPostHog(t)
	o11y.MustStart(t.Context())

	logger := newLogger()
	for _, kind := range []errors.KindType{
		errors.KindUnknown,
		errors.KindInternal,
		errors.KindCritical,
		errors.KindServiceUnavailable,
		errors.KindUnauthenticated,
		errors.KindUnauthorized,
		errors.KindResourceExhausted,
		errors.KindInvalidInput,
		errors.KindNotFound,
		errors.KindConflict,
		errors.KindUnprocessable,
		errors.KindCanceled,
		errors.KindWarn,
	} {
		err := errors.New("boom").WithKind(kind)
		logger.Error(t.Context(), err)
		logger.ErrorBySeverity(t.Context(), err)
	}

	require.NoError(t, o11y.Shutdown(t.Context()))

	exceptions := posthog.exceptions(t)
	kinds := make([]any, 0, len(exceptions))
	for _, exception := range exceptions {
		kinds = append(kinds, exception["kind"])
	}
	require.ElementsMatch(t, []any{
		"UNKNOWN", "UNKNOWN",
		"INTERNAL", "INTERNAL",
		"CRITICAL", "CRITICAL",
		"SERVICE_UNAVAILABLE", "SERVICE_UNAVAILABLE",
	}, kinds)
}

func TestErrorTracking_FatalFlushesBeforePanicking(t *testing.T) {
	posthog := startPostHog(t)
	o11y.MustStart(t.Context())
	t.Cleanup(func() { _ = o11y.Shutdown(context.WithoutCancel(t.Context())) })

	require.Panics(t, func() {
		newLogger().Fatal(t.Context(), errors.New("boom").WithKind(errors.KindFatal))
	})

	exceptions := posthog.exceptions(t)
	require.Len(t, exceptions, 1)
	require.Equal(t, "FATAL", exceptions[0]["kind"])
}

func TestErrorTracking_WithoutAProjectKeyNeverCallsPostHog(t *testing.T) {
	posthog := startPostHog(t)
	t.Setenv("POSTHOG_PROJECT_KEY", "")

	require.NoError(t, o11y.Start(t.Context()))
	newLogger().Error(t.Context(), errors.New("boom").WithKind(errors.KindInternal))
	require.NoError(t, o11y.Shutdown(t.Context()))

	posthog.mu.Lock()
	defer posthog.mu.Unlock()
	require.Zero(t, posthog.requests)
}

func TestErrorTracking_ShutdownStopsWaitingAtTheContextDeadline(t *testing.T) {
	stalled := startStalledPostHog(t)
	o11y.MustStart(t.Context())
	newLogger().Error(t.Context(), errors.New("boom").WithKind(errors.KindInternal))

	ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	require.Error(t, o11y.Shutdown(ctx))
	require.Less(t, time.Since(start), 2*time.Second)
	<-stalled
}

func TestErrorTracking_ShutdownStopsWaitingAfterFiveSeconds(t *testing.T) {
	if testing.Short() {
		t.Skip("waits for the 5s shutdown timeout")
	}

	stalled := startStalledPostHog(t)
	o11y.MustStart(t.Context())
	newLogger().Error(t.Context(), errors.New("boom").WithKind(errors.KindInternal))

	start := time.Now()
	require.Error(t, o11y.Shutdown(t.Context()))
	require.InDelta(t, 5*time.Second, time.Since(start), float64(time.Second))
	<-stalled
}

// startStalledPostHog answers no request until the test ends, and its channel yields once a request arrived.
func startStalledPostHog(t *testing.T) <-chan struct{} {
	t.Helper()

	startPostHog(t)

	arrived := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		once.Do(func() { close(arrived) })
		select {
		case <-release:
		case <-time.After(10 * time.Second):
		}
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	t.Setenv("POSTHOG_HOST", server.URL)

	return arrived
}
