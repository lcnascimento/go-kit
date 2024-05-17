package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/lcnascimento/go-kit/env"
	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/runtime"
)

// Level is a wrapper around slog.Leveler.
// It adds some extra Levels, like Critical.
const (
	LevelDebug    = slog.LevelDebug
	LevelInfo     = slog.LevelInfo
	LevelWarn     = slog.LevelWarn
	LevelError    = slog.LevelError
	LevelCritical = slog.Level(12)
	LevelFatal    = slog.Level(14)
)

const (
	messageKey = "message"
)

// Attrs are wrappers around slog.Attr.
var (
	String = slog.String
	Int    = slog.Int
	Float  = slog.Float64
	Bool   = slog.Bool
	Any    = slog.Any

	errorAttr = func(err error) slog.Attr {
		attrs := []any{
			slog.String("code", string(errors.Code(err))),
			slog.String("kind", string(errors.Kind(err))),
			slog.Bool("retryable", errors.Retryable(err)),
		}

		if root := errors.RootError(err); root != err.Error() {
			attrs = append(attrs, slog.String("root", root))
		}

		if stack := errors.Stack(err); len(stack) > 0 {
			attrs = append(attrs, slog.Any("stack", stackList(stack)))
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
)

var levelNames = map[slog.Leveler]string{
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarn:     "WARN",
	LevelError:    "ERROR",
	LevelCritical: "CRITICAL",
	LevelFatal:    "FATAL",
}

var levelByName = map[string]slog.Level{
	"DEBUG":    LevelDebug,
	"INFO":     LevelInfo,
	"WARN":     LevelWarn,
	"ERROR":    LevelError,
	"CRITICAL": LevelCritical,
	"FATAL":    LevelFatal,
}

type handler struct {
	slog.Handler
}

//nolint:gocritic // we actually must implement this contract.
func (h handler) Handle(ctx context.Context, r slog.Record) error {
	attrs := []any{}
	for k := range *contextKeys {
		if value := ctx.Value(k); value != nil {
			attrs = append(attrs, Any(string(k), value))
		}
	}

	if len(attrs) > 0 {
		r.AddAttrs(slog.Group("context", attrs...))
	}

	return h.Handler.Handle(ctx, r)
}

func newHandler(out io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &handler{
		Handler: slog.NewJSONHandler(out, opts),
	}
}

func initLogLevel() {
	name := env.GetString("LOG_LEVEL", "INFO")

	level := levelByName[strings.ToUpper(name)]
	if level == 0 {
		level = LevelInfo
	}

	logLevel.Set(level)
}

func defaultReplaceAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.MessageKey:
		a.Key = messageKey
	case slog.SourceKey:
		if source, ok := a.Value.Any().(*slog.Source); ok {
			a.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))
		}
	case slog.LevelKey:
		level, _ := a.Value.Any().(slog.Level)
		label, exists := levelNames[level]
		if !exists {
			label = level.String()
		}

		a.Value = slog.StringValue(label)
	}

	return a
}

func stackList(stack []runtime.StackFrame) []string {
	list := []string{}

	for _, s := range stack {
		list = append(list, fmt.Sprintf("%s:%d", s.File, s.LineNumber))
	}

	return list
}
