package errortracking

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/posthog/posthog-go"

	"github.com/lcnascimento/go-kit/env"
	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

const flushTimeout = 5 * time.Second

var (
	mu     sync.Mutex
	client atomic.Pointer[posthog.Client]
)

func Start() error {
	mu.Lock()
	defer mu.Unlock()

	if client.Load() != nil {
		return nil
	}

	key := env.Get[string]("POSTHOG_PROJECT_KEY")
	if key == "" || global.Config().Disabled {
		return nil
	}

	c, err := posthog.NewWithConfig(key, posthog.Config{
		Endpoint: env.Get[string]("POSTHOG_HOST"),
		Logger:   warnCappedLogger{},
	})
	if err != nil {
		return errors.Wrap(err, "failed to create posthog client")
	}

	client.Store(&c)
	return nil
}

func Shutdown(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	c := client.Swap(nil)
	if c == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, flushTimeout)
	defer cancel()

	if err := (*c).CloseWithContext(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown posthog client")
	}

	return nil
}

func Capture(ctx context.Context, err error) {
	c := client.Load()
	if c == nil || !captured(errors.Kind(err)) {
		return
	}

	if err := (*c).Enqueue(newException(ctx, err)); err != nil {
		slog.Default().WarnContext(ctx, "failed to enqueue exception to posthog", "error", err.Error())
	}
}

func Flush(ctx context.Context) {
	c := client.Load()
	if c == nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, flushTimeout)
	defer cancel()

	if err := (*c).FlushWithContext(ctx); err != nil {
		slog.Default().WarnContext(ctx, "failed to flush exceptions to posthog", "error", err.Error())
	}
}

func captured(kind errors.KindType) bool {
	switch kind {
	case errors.KindUnknown, errors.KindInternal, errors.KindCritical, errors.KindServiceUnavailable, errors.KindFatal:
		return true
	default:
		return false
	}
}
