package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/log/format"
	"github.com/lcnascimento/go-kit/propagation"
)

type Formatter interface {
	// Format formats the log payload that will be rendered.
	Format(context.Context, *format.LogInput) any
}

// LoggerInput defines the dependencies of a Logger.
type LoggerInput struct {
	Level      string
	Formatter  Formatter
	Attributes propagation.ContextKeySet

	// This should only be given for testing!
	Now func() time.Time
}

// Logger is the structure responsible for log data.
type Logger struct {
	level      Level
	formatter  Formatter
	attributes propagation.ContextKeySet
	now        func() time.Time
}

// NewLogger constructs a new Logger instance.
func NewLogger(input LoggerInput) *Logger {
	logger := &Logger{
		level:      levelStringValueMap[input.Level],
		attributes: input.Attributes,
		formatter:  input.Formatter,
		now:        time.Now,
	}

	if logger.level < LevelCritical || logger.level > LevelDebug {
		logger.level = LevelInfo
	}

	if logger.formatter == nil {
		logger.formatter = format.NewDefault()
	}

	if input.Now != nil {
		logger.now = input.Now
	}

	return logger
}

// Debug logs debug data.
func (l Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelDebug {
		l.printMsg(ctx, fmt.Sprintf(msg, args...), LevelDebug)
	}
}

// Info logs info data.
func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelInfo {
		l.printMsg(ctx, fmt.Sprintf(msg, args...), LevelInfo)
	}
}

// Warning logs warning data.
func (l Logger) Warning(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelWarning {
		l.printMsg(ctx, fmt.Sprintf(msg, args...), LevelWarning)
	}
}

// Error logs error data.
func (l Logger) Error(ctx context.Context, err error) {
	if l.level >= LevelError {
		l.printError(ctx, err, LevelError)
	}
}

// Critical logs critical data.
func (l Logger) Critical(ctx context.Context, err error) {
	if l.level >= LevelCritical {
		l.printError(ctx, err, LevelCritical)
	}
}

// Fatal logs critical data and exists current program execution.
func (l Logger) Fatal(ctx context.Context, err error) {
	if l.level >= LevelCritical {
		l.printError(ctx, err, LevelCritical)
		os.Exit(1)
	}
}

// JSON logs JSON data in Debug level.
func (l Logger) JSON(ctx context.Context, data any) {
	if l.level < LevelDebug {
		return
	}

	if _, err := json.Marshal(data); err != nil {
		l.Error(ctx, errors.New("could not marshal payload to JSON format").WithRootError(err))
		return
	}

	l.printJSON(ctx, data)
}

func (l Logger) printMsg(ctx context.Context, msg string, level Level) {
	payload := l.formatter.Format(ctx, &format.LogInput{
		Level:      level.String(),
		Message:    msg,
		Attributes: l.attributes,
		Timestamp:  l.now(),
	})

	data, _ := json.Marshal(payload)
	fmt.Println(string(data))
}

func (l Logger) printJSON(ctx context.Context, jsonData any) {
	payload := l.formatter.Format(ctx, &format.LogInput{
		Level:      LevelDebug.String(),
		Message:    "JSON data logged",
		Payload:    jsonData,
		Attributes: l.attributes,
		Timestamp:  l.now(),
	})

	data, _ := json.Marshal(payload)
	fmt.Println(string(data))
}

func (l Logger) printError(ctx context.Context, err error, level Level) {
	payload := l.formatter.Format(ctx, &format.LogInput{
		Level:      level.String(),
		Message:    err.Error(),
		Err:        err,
		Attributes: l.attributes,
		Timestamp:  l.now(),
	})

	data, _ := json.Marshal(payload)
	fmt.Println(string(data))
}
