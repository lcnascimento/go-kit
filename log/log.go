package log

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-slog/otelslog"

	"github.com/lcnascimento/go-kit/propagation"
)

var (
	logLevel    = &slog.LevelVar{}
	contextKeys = &propagation.ContextKeySet{}
)

func init() {
	initLogLevel()

	options := &slog.HandlerOptions{
		AddSource:   true,
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

// SetContextKeySet changes the logger context key set.
func SetContextKeySet(set propagation.ContextKeySet) {
	if set == nil {
		return
	}

	*contextKeys = set
}

// Debug logs debug data.
func Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, LevelDebug, msg, attributes(attrs...))
}

// Info logs info data.
func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, LevelInfo, msg, attributes(attrs...))
}

// Warn logs warning data.
func Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, LevelWarn, msg, attributes(attrs...))
}

// Error logs error data based on an error object.
func Error(ctx context.Context, err error, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, LevelError, err.Error(), attributes(attrs...), errorAttr(err))
}

// Errorw logs error data based on an error message.
func Errorw(ctx context.Context, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, LevelError, msg, attributes(attrs...))
}

// Critical logs critical data based on an error object.
func Critical(ctx context.Context, err error, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, LevelCritical, err.Error(), attributes(attrs...), errorAttr(err))
}

// Criticalw logs critical data based on an error message.
func Criticalw(ctx context.Context, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, LevelCritical, msg, attributes(attrs...))
}
