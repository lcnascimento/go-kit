package log

import (
	"context"
	"errors"
	"log/slog"
	r "runtime"

	"go.opentelemetry.io/contrib/bridges/otelslog"

	"github.com/lcnascimento/go-kit/runtime"
)

// Logger is a wrapper around slog.Logger that adds OpenTelemetry features.
type Logger struct {
	level  slog.Level
	logger *slog.Logger
}

// NewLogger creates a new Logger.
func NewLogger(pkg string) *Logger {
	return &Logger{
		level:  getLevel(),
		logger: otelslog.NewLogger(pkg),
	}
}

// Debug logs debug data.
func (l *Logger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelDebug {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelDebug, msg, source, baggage, attributes)
}

// Info logs info data.
func (l *Logger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelInfo {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelInfo, msg, source, baggage, attributes)
}

// Warn logs warning data.
func (l *Logger) Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelWarn {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelWarn, msg, source, baggage, attributes)
}

// Error logs error data based on an error object.
func (l *Logger) Error(ctx context.Context, err error, attrs ...slog.Attr) {
	if l.level > LevelError {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	errAttr := errorAttr(ctx, LevelError, err)
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelError, err.Error(), source, baggage, attributes, errAttr)
}

// Errorw logs error data based on an error message.
func (l *Logger) Errorw(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelError {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	errAttr := errorAttr(ctx, LevelCritical, errors.New(msg))
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelError, msg, source, baggage, attributes, errAttr)
}

// Critical logs critical data based on an error object.
func (l *Logger) Critical(ctx context.Context, err error, attrs ...slog.Attr) {
	if l.level > LevelCritical {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	errAttr := errorAttr(ctx, LevelCritical, err)
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelCritical, err.Error(), source, baggage, attributes, errAttr)
}

// Criticalw logs critical data based on an error message.
func (l *Logger) Criticalw(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelCritical {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	errAttr := errorAttr(ctx, LevelCritical, errors.New(msg))
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelCritical, msg, source, baggage, attributes, errAttr)
}

// Fatal logs fatal data based on an error object.
// It terminates the current goroutine.
//
// DANGER: This function uses `runtime.Goexit()`. Therefore, if called within the main goroutine
// you MUST call `defer os.Exit(0)` at the top of the main function.
func (l *Logger) Fatal(ctx context.Context, err error, attrs ...slog.Attr) {
	if l.level > LevelFatal {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	errAttr := errorAttr(ctx, LevelFatal, err)
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelFatal, err.Error(), source, baggage, attributes, errAttr)

	r.Goexit()
}

// Fatalw logs fatal data based on an error message.
// It terminates the current goroutine.
//
// DANGER: This function uses `runtime.Goexit()`. Therefore, if called within the main goroutine
// you MUST call `defer os.Exit(0)` at the top of the main function.
func (l *Logger) Fatalw(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelFatal {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	errAttr := errorAttr(ctx, LevelFatal, errors.New(msg))
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelFatal, msg, source, baggage, attributes, errAttr)

	r.Goexit()
}
