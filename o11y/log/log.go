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

// Trace should be used for very detailed, low-level information, usually about the internal state
// of the application or detailed execution paths. This level is often used for tracing the flow of the application,
// and it's usually the most verbose.
func (l *Logger) Trace(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelTrace {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelTrace, msg, source, baggage, attributes)
}

// Debug provides information that is useful for debugging but less detailed than trace.
// It's used for diagnostic information that might be helpful when trying to understand the flow
// of an application or troubleshoot issues.
func (l *Logger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelDebug {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelDebug, msg, source, baggage, attributes)
}

// Info gives general information about the application’s state or normal operations,
// typically showing what's happening at a higher level.
func (l *Logger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelInfo {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelInfo, msg, source, baggage, attributes)
}

// Warn indicates potentially harmful situations, but not necessarily errors.
// Used for things like deprecated features or unusual but non-critical conditions.
func (l *Logger) Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l.level > LevelWarn {
		return
	}

	source := slog.String("source", runtime.Caller().String())
	baggage := baggageAttr(ctx)
	attributes := attributes(attrs...)

	l.logger.LogAttrs(ctx, LevelWarn, msg, source, baggage, attributes)
}

// Error indicates significant issues that prevent the application from functioning properly.
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

// Errorw is semantically equivalent to [Error], but it logs an error message instead of an error object.
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

// Critical indicates critical errors that may lead to system failure or data loss.
// This level is used for situations that are extremely serious and require immediate attention.
// It does not lead to a crash, but it should be logged and monitored closely.
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

// Criticalw is semantically equivalent to [Critical], but it logs an error message instead of an error object.
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

// Fatal indicates a fatal error that causes the current goroutine to terminate.
// This level is used for situations that are extremely serious and require immediate attention.
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

// Fatalw is semantically equivalent to [Fatal], but it logs an error message instead of an error object.
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
