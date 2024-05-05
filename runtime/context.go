package runtime

import (
	"context"
	"os"
	"os/signal"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/lcnascimento/go-kit/env"
	"github.com/lcnascimento/go-kit/propagation"
)

var (
	serviceName    = env.GetString("SERVICE_NAME", "unknown")
	serviceVersion = env.GetString("SERVICE_VERSION", "v0.0.0")
)

// SystemContextKeySet returns base keys that identifies the current service.
func SystemContextKeySet() propagation.ContextKeySet {
	return propagation.ContextKeySet{
		propagation.ContextKey(semconv.ServiceNameKey):    serviceName,
		propagation.ContextKey(semconv.ServiceVersionKey): serviceVersion,
	}
}

// ContextWithOSSignalCancellation builds a new Context that cancels itself on os.Interrupt signals.
func ContextWithOSSignalCancellation() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		cancel()
	}()

	return ctx
}
