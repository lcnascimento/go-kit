package temporalclient

import (
	"go.opentelemetry.io/otel/trace/noop"
	"go.temporal.io/sdk/interceptor"

	"github.com/lcnascimento/go-kit/temporal/internal/interceptors"
)

type Option func(*options)

type options struct {
	splitTracer bool
	noopMetrics bool
}

func WithSplitTracer() Option {
	return func(o *options) {
		o.splitTracer = true
	}
}

func WithNoopMetrics() Option {
	return func(o *options) {
		o.noopMetrics = true
	}
}

func Interceptors(opts ...Option) []interceptor.ClientInterceptor {
	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	if o.splitTracer {
		return []interceptor.ClientInterceptor{
			interceptors.NewClientTracer(noop.NewTracerProvider()),
		}
	}

	return []interceptor.ClientInterceptor{
		interceptors.NewTracer(noop.NewTracerProvider()),
	}
}
