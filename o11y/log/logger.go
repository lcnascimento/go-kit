package log

import (
	"context"
	"log/slog"

	"github.com/lcnascimento/go-kit/errors"
)

// Logger is logging structure with a more rigid interface than slog.
// It forces the use of context and attributes, in a more opinionated style.
// It's usage is highly recommended due to its facilities on error handling.
type Logger interface {
	// Debug logs a message at the debug level.
	Debug(ctx context.Context, msg string, attrs ...Attr)

	// Info logs a message at the info level.
	Info(ctx context.Context, msg string, attrs ...Attr)

	// Warn logs a message at the warn level.
	Warn(ctx context.Context, msg string, attrs ...Attr)

	// ErrorMessage logs a message at the error level.
	ErrorMessage(ctx context.Context, msg string, attrs ...Attr)

	// Error logs error data at the error level.
	Error(ctx context.Context, err error, attrs ...Attr)

	// ErrorBySeverity logs error data with the level based on the severity of the error.
	ErrorBySeverity(ctx context.Context, err error, attrs ...Attr)

	// ErrorAttr wraps an error with error key.
	ErrorAttr(err error) Attr

	// CriticalMessage logs a message at the custom critical level.
	CriticalMessage(ctx context.Context, msg string, attrs ...Attr)

	// Critical logs error data at the custom critical level.
	Critical(ctx context.Context, err error, attrs ...Attr)

	// FatalMessage logs a message at the custom fatal level and panics.
	FatalMessage(ctx context.Context, msg string, attrs ...Attr)

	// Fatal logs error data at the custom fatal level and panics.
	Fatal(ctx context.Context, err error, attrs ...Attr)
}

// NewLogger returns a new [slog.Logger] backed by a new [Handler]. See [NewHandler] for details on how the backing Handler is created.
func NewLogger(name string, opts ...HandlerOption) (Logger, error) {
	handler, err := NewHandler(name, opts...)
	if err != nil {
		return nil, err
	}

	return &implementation{logger: slog.New(handler)}, nil
}

func MustNewLogger(name string, opts ...HandlerOption) Logger {
	logger, err := NewLogger(name, opts...)
	if err != nil {
		panic(errors.Wrap(err, "failed to create logger"))
	}

	return logger
}

type implementation struct {
	logger *slog.Logger
}

// Debug logs a message at the debug level.
func (l *implementation) Debug(ctx context.Context, msg string, attrs ...Attr) {
	l.logger.LogAttrs(ctx, LevelDebug, msg, attrs...)
}

// Info logs a message at the info level.
func (l *implementation) Info(ctx context.Context, msg string, attrs ...Attr) {
	l.logger.LogAttrs(ctx, LevelInfo, msg, attrs...)
}

// Warn logs a message at the warn level.
func (l *implementation) Warn(ctx context.Context, msg string, attrs ...Attr) {
	l.logger.LogAttrs(ctx, LevelWarn, msg, attrs...)
}

// ErrorMessage logs a message at the error level.
func (l *implementation) ErrorMessage(ctx context.Context, msg string, attrs ...Attr) {
	err := errors.New("%s", msg)
	l.Error(ctx, err, attrs...)
}

// Error logs error data at the error level.
func (l *implementation) Error(ctx context.Context, err error, attrs ...Attr) {
	if err == nil {
		return
	}

	attrs = append(attrs, errorAttr(err))
	l.logger.LogAttrs(ctx, LevelError, err.Error(), attrs...)
}

// ErrorAttr wraps an error with error key.
func (l *implementation) ErrorAttr(err error) Attr {
	return Any(ErrorKey, err)
}

// CriticalMessage logs a message at the custom critical level.
func (l *implementation) CriticalMessage(ctx context.Context, msg string, attrs ...Attr) {
	err := errors.New("%s", msg).WithKind(errors.KindCritical)
	l.Critical(ctx, err, attrs...)
}

// Critical logs error data at the custom critical level.
func (l *implementation) Critical(ctx context.Context, err error, attrs ...Attr) {
	if err == nil {
		return
	}

	attrs = append(attrs, errorAttr(err))
	l.logger.LogAttrs(ctx, LevelCritical, err.Error(), attrs...)
}

// FatalMessage logs a message at the custom fatal level and panics.
func (l *implementation) FatalMessage(ctx context.Context, msg string, attrs ...Attr) {
	err := errors.New("%s", msg).WithKind(errors.KindFatal)
	l.Fatal(ctx, err, attrs...)
}

// Fatal logs error data at the custom fatal level and panics.
func (l *implementation) Fatal(ctx context.Context, err error, attrs ...Attr) {
	if err == nil {
		return
	}

	attrs = append(attrs, errorAttr(err))
	l.logger.LogAttrs(ctx, LevelFatal, err.Error(), attrs...)

	panic(err.Error())
}

// ErrorBySeverity logs error data with the level based on the severity of the error.
func (l *implementation) ErrorBySeverity(ctx context.Context, err error, attrs ...Attr) {
	switch errors.Severity(err) {
	case errors.SeverityWarn:
		attrs = append(attrs, errorAttr(err))
		l.Warn(ctx, err.Error(), attrs...)
	case errors.SeverityCritical:
		l.Critical(ctx, err, attrs...)
	case errors.SeverityFatal:
		l.Fatal(ctx, err, attrs...)
	default:
		l.Error(ctx, err, attrs...)
	}
}

func errorAttr(err error) Attr {
	attrs := []any{
		String("message", err.Error()),
		String("kind", string(errors.Code(err))),
		String("avn_code", string(errors.Code(err))),
		String("avn_kind", string(errors.Kind(err))),
		String("severity", errors.Severity(err).String()),
		Bool("retryable", errors.IsRetryable(err)),
	}

	for key, value := range errors.Attributes(err) {
		attrs = append(attrs, String(key, value))
	}

	reasons := errors.Reasons(err)
	if len(reasons) > 0 {
		attrs = append(attrs, Any("reasons", reasons))
	}

	return Group("error", attrs...)
}
