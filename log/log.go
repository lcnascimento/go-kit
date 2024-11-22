package log

import (
	"context"
	"log/slog"
	"os"
	r "runtime"

	"github.com/go-slog/otelslog"

	"github.com/lcnascimento/go-kit/runtime"
)

var logLevel = &slog.LevelVar{}

func init() {
	initLogLevel()

	options := &slog.HandlerOptions{
		AddSource:   false,
		Level:       logLevel,
		ReplaceAttr: defaultReplaceAttr,
	}

	handler := newHandler(os.Stdout, options)
	handler = otelslog.NewHandler(handler)

	slog.SetDefault(slog.New(handler))
}

// SetLevel changes the level of the logger.
func SetLevel(level slog.Level) {
	logLevel.Set(level)
}

// Debug logs debug data.
func Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	source := slog.String("source", runtime.Caller().String())
	attributes := attributes(attrs...)

	slog.LogAttrs(ctx, LevelDebug, msg, source, attributes)
}

// Info logs info data.
func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	source := slog.String("source", runtime.Caller().String())
	attributes := attributes(attrs...)

	slog.LogAttrs(ctx, LevelInfo, msg, source, attributes)
}

// Warn logs warning data.
func Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	source := slog.String("source", runtime.Caller().String())
	attributes := attributes(attrs...)

	slog.LogAttrs(ctx, LevelWarn, msg, source, attributes)
}

// Error logs error data based on an error object.
func Error(ctx context.Context, err error, attrs ...slog.Attr) {
	source := slog.String("source", runtime.Caller().String())
	errAttr := errorAttr(err)
	attributes := attributes(attrs...)

	slog.LogAttrs(ctx, LevelError, err.Error(), source, errAttr, attributes)
}

// Errorw logs error data based on an error message.
func Errorw(ctx context.Context, msg string, attrs ...slog.Attr) {
	source := slog.String("source", runtime.Caller().String())
	attributes := attributes(attrs...)

	slog.LogAttrs(ctx, LevelError, msg, source, attributes)
}

// Critical logs critical data based on an error object.
func Critical(ctx context.Context, err error, attrs ...slog.Attr) {
	source := slog.String("source", runtime.Caller().String())
	errAttr := errorAttr(err)
	attributes := attributes(attrs...)

	slog.LogAttrs(ctx, LevelCritical, err.Error(), source, errAttr, attributes)
}

// Criticalw logs critical data based on an error message.
func Criticalw(ctx context.Context, msg string, attrs ...slog.Attr) {
	source := slog.String("source", runtime.Caller().String())
	attributes := attributes(attrs...)

	slog.LogAttrs(ctx, LevelCritical, msg, source, attributes)
}

// Fatal logs fatal data based on an error object.
// It terminates the current goroutine.
//
// DANGER: This function uses `runtime.Goexit()`. Therefore, if called within the main goroutine
// you MUST call `defer os.Exit(0)` at the top of the main function.
func Fatal(ctx context.Context, err error, attrs ...slog.Attr) {
	source := slog.String("source", runtime.Caller().String())
	errAttr := errorAttr(err)
	attributes := attributes(attrs...)

	slog.LogAttrs(ctx, LevelFatal, err.Error(), source, attributes, errAttr)

	r.Goexit()
}

// Fatalw logs fatal data based on an error message.
// It terminates the current goroutine.
//
// DANGER: This function uses `runtime.Goexit()`. Therefore, if called within the main goroutine
// you MUST call `defer os.Exit(0)` at the top of the main function.
func Fatalw(ctx context.Context, msg string, attrs ...slog.Attr) {
	source := slog.String("source", runtime.Caller().String())
	attributes := attributes(attrs...)

	slog.LogAttrs(ctx, LevelFatal, msg, source, attributes)

	r.Goexit()
}
