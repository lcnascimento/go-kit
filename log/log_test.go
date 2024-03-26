package log_test

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/log"
	"github.com/lcnascimento/go-kit/propagation"
)

func TestNewLogger(t *testing.T) {
	t.Run("should use LogLevel INFO when not specified", func(t *testing.T) {
		ctx := context.Background()

		logger := log.NewLogger(log.LoggerInput{Now: mockNow})

		t.Run("should not log Debug", func(t *testing.T) {
			out := captureOutput(func() {
				logger.Debug(ctx, "random message")
			})

			assert.Empty(t, out)
		})

		t.Run("should log Info", func(t *testing.T) {
			out := captureOutput(func() {
				logger.Info(ctx, "random message")
			})

			assert.Equal(t, `{"level":"INFO","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`, out)
		})
	})
}

func TestDebug(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       propagation.ContextKeySet
		msg         string
		msgArgs     []any
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message",
			expectedLog: `{"level":"DEBUG","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message with dynamic data %d",
			msgArgs:     []any{1},
			expectedLog: `{"level":"DEBUG","message":"random message with dynamic data 1","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, propagation.ContextKey("attr1"), "value1"),
			level:       "DEBUG",
			msg:         "random message",
			attrs:       propagation.ContextKeySet{"attr1": true},
			expectedLog: `{"attributes":{"attr1":"value1"},"level":"DEBUG","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{
				Level:      tc.level,
				Attributes: tc.attrs,
				Now:        mockNow,
			})

			out := captureOutput(func() {
				logger.Debug(tc.ctx, tc.msg, tc.msgArgs...)
			})

			assert.Equal(t, tc.expectedLog, out)
		})
	}
}

func TestInfo(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       propagation.ContextKeySet
		msg         string
		msgArgs     []any
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message",
			expectedLog: `{"level":"INFO","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			msg:         "random message",
			expectedLog: `{"level":"INFO","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message with dynamic data %d",
			msgArgs:     []any{1},
			expectedLog: `{"level":"INFO","message":"random message with dynamic data 1","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, propagation.ContextKey("attr1"), "value1"),
			level:       "DEBUG",
			msg:         "random message",
			attrs:       propagation.ContextKeySet{propagation.ContextKey("attr1"): true},
			expectedLog: `{"attributes":{"attr1":"value1"},"level":"INFO","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{
				Level:      tc.level,
				Attributes: tc.attrs,
				Now:        mockNow,
			})

			out := captureOutput(func() {
				logger.Info(tc.ctx, tc.msg, tc.msgArgs...)
			})

			assert.Equal(t, tc.expectedLog, out)
		})
	}
}

func TestWarning(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       propagation.ContextKeySet
		msg         string
		msgArgs     []any
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message",
			expectedLog: `{"level":"WARNING","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			msg:         "random message",
			expectedLog: `{"level":"WARNING","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			msg:         "random message",
			expectedLog: `{"level":"WARNING","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message with dynamic data %d",
			msgArgs:     []any{1},
			expectedLog: `{"level":"WARNING","message":"random message with dynamic data 1","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, propagation.ContextKey("attr1"), "value1"),
			level:       "DEBUG",
			msg:         "random message",
			attrs:       propagation.ContextKeySet{propagation.ContextKey("attr1"): true},
			expectedLog: `{"attributes":{"attr1":"value1"},"level":"WARNING","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{
				Level:      tc.level,
				Attributes: tc.attrs,
				Now:        mockNow,
			})

			out := captureOutput(func() {
				logger.Warning(tc.ctx, tc.msg, tc.msgArgs...)
			})

			assert.Equal(t, tc.expectedLog, out)
		})
	}
}

func TestError(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       propagation.ContextKeySet
		err         error
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			err:         errors.New("random error"),
			expectedLog: "",
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, propagation.ContextKey("attr1"), "value1"),
			level:       "DEBUG",
			err:         errors.New("random error"),
			attrs:       propagation.ContextKeySet{propagation.ContextKey("attr1"): true},
			expectedLog: `{"attributes":{"attr1":"value1","error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{
				Level:      tc.level,
				Attributes: tc.attrs,
				Now:        mockNow,
			})

			out := captureOutput(func() {
				logger.Error(tc.ctx, tc.err)
			})

			assert.Equal(t, tc.expectedLog, out)
		})
	}
}

func TestCritical(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       propagation.ContextKeySet
		err         error
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, propagation.ContextKey("attr1"), "value1"),
			level:       "DEBUG",
			err:         errors.New("random error"),
			attrs:       propagation.ContextKeySet{propagation.ContextKey("attr1"): true},
			expectedLog: `{"attributes":{"attr1":"value1","error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{
				Level:      tc.level,
				Attributes: tc.attrs,
				Now:        mockNow,
			})

			out := captureOutput(func() {
				logger.Critical(tc.ctx, tc.err)
			})

			assert.Equal(t, tc.expectedLog, out)
		})
	}
}

func TestJSON(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		sysLogLevel string
		level       []log.Level
		attrs       propagation.ContextKeySet
		data        any
		expectedLog string
	}{
		{
			desc:        "should log when system LogLevel is DEBUG",
			ctx:         ctx,
			sysLogLevel: "DEBUG",
			data:        map[string]string{"foo": "bar"},
			expectedLog: `{"level":"DEBUG","message":"JSON data logged","payload":{"foo":"bar"},"timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log an error when data can not be JSON marshalled",
			ctx:         ctx,
			sysLogLevel: "DEBUG",
			data:        make(chan int),
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"json: unsupported type: chan int"},"level":"ERROR","message":"could not marshal payload to JSON format","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log anything when system LogLevel is INFO",
			ctx:         ctx,
			sysLogLevel: "INFO",
			data:        map[string]string{"foo": "bar"},
			expectedLog: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{
				Level:      tc.sysLogLevel,
				Attributes: tc.attrs,
				Now:        mockNow,
			})

			out := captureOutput(func() {
				logger.JSON(tc.ctx, tc.data)
			})

			assert.Equal(t, tc.expectedLog, out)
		})
	}
}

func captureOutput(output func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	output()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	return strings.TrimRight(string(out), "\n")
}

func mockNow() time.Time {
	return time.Date(2021, 0, 1, 12, 0, 0, 0, time.UTC)
}
