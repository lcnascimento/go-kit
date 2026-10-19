package log_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/lcnascimento/go-kit/errors"

	"github.com/lcnascimento/go-kit/o11y/log"
)

func TestLogger_ErrorLevelsRecordTheErrorOnTheSpan(t *testing.T) {
	t.Parallel()

	logger := log.MustNewLogger("test", log.WithLogHandler(slog.NewTextHandler(io.Discard, nil)))

	cases := map[string]func(logger log.Logger, ctx context.Context){
		"Error":                 func(l log.Logger, ctx context.Context) { l.Error(ctx, errors.New("boom")) },
		"ErrorMessage":          func(l log.Logger, ctx context.Context) { l.ErrorMessage(ctx, "boom") },
		"Critical":              func(l log.Logger, ctx context.Context) { l.Critical(ctx, errors.New("boom")) },
		"CriticalMessage":       func(l log.Logger, ctx context.Context) { l.CriticalMessage(ctx, "boom") },
		"ErrorBySeverity error": func(l log.Logger, ctx context.Context) { l.ErrorBySeverity(ctx, errors.New("boom")) },
		"ErrorBySeverity critical": func(l log.Logger, ctx context.Context) {
			l.ErrorBySeverity(ctx, errors.New("boom").WithKind(errors.KindCritical))
		},
	}

	for name, logError := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			recorder := tracetest.NewSpanRecorder()
			ctx, span := trace.NewTracerProvider(trace.WithSpanProcessor(recorder)).Tracer("test").Start(t.Context(), "op")

			logError(logger, ctx)
			span.End()

			ended := recorder.Ended()[0]
			require.Equal(t, codes.Error, ended.Status().Code)
			require.Equal(t, "boom", ended.Status().Description)
			require.Len(t, ended.Events(), 1)
			require.Equal(t, "exception", ended.Events()[0].Name)
		})
	}
}

func TestLogger_WarnSeverityLeavesTheSpanStatusUnset(t *testing.T) {
	t.Parallel()

	logger := log.MustNewLogger("test", log.WithLogHandler(slog.NewTextHandler(io.Discard, nil)))
	recorder := tracetest.NewSpanRecorder()
	ctx, span := trace.NewTracerProvider(trace.WithSpanProcessor(recorder)).Tracer("test").Start(t.Context(), "op")

	logger.ErrorBySeverity(ctx, errors.New("not found").WithKind(errors.KindNotFound))
	span.End()

	ended := recorder.Ended()[0]
	require.Equal(t, codes.Unset, ended.Status().Code)
	require.Empty(t, ended.Events())
}
