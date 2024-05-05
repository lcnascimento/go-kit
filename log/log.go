package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/log/format"
	"github.com/lcnascimento/go-kit/propagation"
	"github.com/lcnascimento/go-kit/runtime"
)

// Formatter defines how structures that formats logs should behavior.
type Formatter interface {
	// Format formats the log payload that will be rendered.
	Format(context.Context, *format.LogInput) any
}

// Logger is the structure responsible for log data.
type Logger struct {
	level       Level
	formatter   Formatter
	contextKeys propagation.ContextKeySet
	now         func() time.Time
}

// NewLogger constructs a new Logger instance.
func NewLogger(opts ...Option) *Logger {
	logger := &Logger{
		now: time.Now,
	}

	for _, opt := range opts {
		opt(logger)
	}

	if logger.level < LevelCritical || logger.level > LevelDebug {
		logger.level = LevelInfo
	}

	if logger.formatter == nil {
		logger.formatter = format.NewDefault()
	}

	if logger.contextKeys == nil {
		logger.contextKeys = propagation.ContextKeySet{}
	}

	return logger
}

// Debug logs debug data.
func (l Logger) Debug(ctx context.Context, msg string, args ...any) {
	if l.level >= LevelDebug {
		l.print(ctx, LevelDebug, msg, args...)
	}
}

// Info logs info data.
func (l Logger) Info(ctx context.Context, msg string, args ...any) {
	if l.level >= LevelInfo {
		l.print(ctx, LevelInfo, msg, args...)
	}
}

// Warning logs warning data.
func (l Logger) Warning(ctx context.Context, msg string, args ...any) {
	if l.level >= LevelWarning {
		l.print(ctx, LevelWarning, msg, args...)
	}
}

// Error logs error data.
func (l Logger) Error(ctx context.Context, err error, args ...any) {
	if l.level >= LevelError {
		l.printError(ctx, LevelError, err, args...)
	}
}

// Critical logs critical data.
func (l Logger) Critical(ctx context.Context, err error, args ...any) {
	if l.level >= LevelCritical {
		l.printError(ctx, LevelCritical, err, args...)
	}
}

// Fatal logs critical data and exists current program execution.
func (l Logger) Fatal(ctx context.Context, err error, args ...any) {
	if l.level >= LevelCritical {
		l.printError(ctx, LevelCritical, err, args...)
		os.Exit(1)
	}
}

func (l Logger) print(ctx context.Context, level Level, msg string, args ...any) {
	msg, attrs := buildMsgAndAttributes(msg, args...)

	payload := l.formatter.Format(ctx, &format.LogInput{
		Level:       level.String(),
		Message:     msg,
		ContextKeys: l.contextKeys,
		Attributes:  attrs,
		Timestamp:   l.now(),
	})

	data, _ := json.Marshal(payload)
	fmt.Println(string(data))
}

func (l Logger) printError(ctx context.Context, level Level, err error, args ...any) {
	msg, attrs := buildMsgAndAttributes(err.Error(), args...)

	attrs["error.root"] = errors.RootError(err)
	attrs["error.kind"] = string(errors.Kind(err))
	attrs["error.code"] = string(errors.Code(err))
	attrs["error.retryable"] = strconv.FormatBool(errors.Retryable(err))

	stack := errors.Stack(err)
	if len(stack) > 0 {
		attrs["error.stack"] = stackList(stack)
	}

	l.print(ctx, level, msg, attrs)
}

func buildMsgAndAttributes(msg string, args ...any) (string, format.AttributeSet) {
	if len(args) == 0 {
		return msg, format.AttributeSet{}
	}

	if attrs, ok := args[0].(format.AttributeSet); ok {
		return msg, attrs
	}

	return fmt.Sprintf(msg, args...), format.AttributeSet{}
}

func stackList(stack []runtime.StackFrame) []string {
	list := []string{}

	for _, s := range stack {
		list = append(list, fmt.Sprintf("%s:%d (%s)", s.File, s.LineNumber, s.Name))
	}

	return list
}
