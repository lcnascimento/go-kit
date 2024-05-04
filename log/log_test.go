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

func TestDebug(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		contextKeys propagation.ContextKeySet
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
			desc:        "should log with context keys",
			ctx:         context.WithValue(ctx, propagation.ContextKey("key1"), "value1"),
			level:       "DEBUG",
			msg:         "random message",
			contextKeys: propagation.ContextKeySet{"key1": true},
			expectedLog: `{"context":{"key1":"value1"},"level":"DEBUG","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(
				log.WithLevel(tc.level),
				log.WithContextKeySet(tc.contextKeys),
				log.WithTimmer(mockNow),
			)

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
		contextKeys propagation.ContextKeySet
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
			desc:        "should log with context keys",
			ctx:         context.WithValue(ctx, propagation.ContextKey("key1"), "value1"),
			level:       "DEBUG",
			msg:         "random message",
			contextKeys: propagation.ContextKeySet{propagation.ContextKey("key1"): true},
			expectedLog: `{"context":{"key1":"value1"},"level":"INFO","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(
				log.WithLevel(tc.level),
				log.WithContextKeySet(tc.contextKeys),
				log.WithTimmer(mockNow),
			)

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
		contextKeys propagation.ContextKeySet
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
			desc:        "should log with context keys",
			ctx:         context.WithValue(ctx, propagation.ContextKey("key1"), "value1"),
			level:       "DEBUG",
			msg:         "random message",
			contextKeys: propagation.ContextKeySet{propagation.ContextKey("key1"): true},
			expectedLog: `{"context":{"key1":"value1"},"level":"WARNING","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(
				log.WithLevel(tc.level),
				log.WithContextKeySet(tc.contextKeys),
				log.WithTimmer(mockNow),
			)

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
		contextKeys propagation.ContextKeySet
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
			desc:        "should log with context keys",
			ctx:         context.WithValue(ctx, propagation.ContextKey("key1"), "value1"),
			level:       "DEBUG",
			err:         errors.New("random error"),
			contextKeys: propagation.ContextKeySet{propagation.ContextKey("key1"): true},
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"context":{"key1":"value1"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(
				log.WithLevel(tc.level),
				log.WithContextKeySet(tc.contextKeys),
				log.WithTimmer(mockNow),
			)

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
		contextKeys propagation.ContextKeySet
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
			desc:        "should log with context keys",
			ctx:         context.WithValue(ctx, propagation.ContextKey("key1"), "value1"),
			level:       "DEBUG",
			err:         errors.New("random error"),
			contextKeys: propagation.ContextKeySet{propagation.ContextKey("key1"): true},
			expectedLog: `{"attributes":{"error.code":"UNKNOWN","error.kind":"UNEXPECTED","error.retryable":"false","error.root":"random error"},"context":{"key1":"value1"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(
				log.WithLevel(tc.level),
				log.WithContextKeySet(tc.contextKeys),
				log.WithTimmer(mockNow),
			)

			out := captureOutput(func() {
				logger.Critical(tc.ctx, tc.err)
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
