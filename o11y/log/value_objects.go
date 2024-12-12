package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/baggage"

	"github.com/lcnascimento/go-kit/errors"
)

// Attrs are wrappers around slog.Attr.
var (
	String = slog.String
	Int    = slog.Int
	Float  = slog.Float64
	Bool   = slog.Bool
	Any    = slog.Any
)

// Level is a wrapper around slog.Leveler.
// It adds some extra Levels, like Critical.
const (
	LevelTrace    = slog.Level(-8)
	LevelDebug    = slog.LevelDebug
	LevelInfo     = slog.LevelInfo
	LevelWarn     = slog.LevelWarn
	LevelError    = slog.LevelError
	LevelCritical = slog.Level(12)
	LevelFatal    = slog.Level(14)
)

var levelNames = map[slog.Leveler]string{
	LevelTrace:    "TRACE",
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarn:     "WARN",
	LevelError:    "ERROR",
	LevelCritical: "CRITICAL",
	LevelFatal:    "FATAL",
}

var levelByName = map[string]slog.Level{
	"TRACE":    LevelTrace,
	"DEBUG":    LevelDebug,
	"INFO":     LevelInfo,
	"WARN":     LevelWarn,
	"ERROR":    LevelError,
	"CRITICAL": LevelCritical,
	"FATAL":    LevelFatal,
}

// utility functions.
var (
	errorAttr = func(ctx context.Context, sev slog.Level, err error) slog.Attr {
		onError(ctx, sev, err)

		attrs := []any{
			slog.String("code", string(errors.Code(err))),
			slog.String("kind", string(errors.Kind(err))),
			slog.Bool("retryable", errors.IsRetryable(err)),
		}

		return slog.Group("error", attrs...)
	}

	attributes = func(attrs ...slog.Attr) slog.Attr {
		iAttrs := []any{}
		for _, attr := range attrs {
			iAttrs = append(iAttrs, attr)
		}

		return slog.Group("attributes", iAttrs...)
	}

	baggageAttr = func(ctx context.Context) slog.Attr {
		bag := baggage.FromContext(ctx)

		attrs := []any{}

		for _, member := range bag.Members() {
			attrs = append(attrs, slog.String(member.Key(), member.Value()))
		}

		if len(attrs) == 0 {
			return slog.Attr{}
		}

		return slog.Group("baggage", attrs...)
	}
)
