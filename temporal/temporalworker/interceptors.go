package temporalworker

import (
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.temporal.io/sdk/interceptor"

	"github.com/lcnascimento/go-kit/temporal/internal/interceptors"
)

type Option func(*options)

type options struct {
	tracerProvider trace.TracerProvider
	splitTracer    bool
}

func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(o *options) {
		o.tracerProvider = tp
	}
}

func WithSplitTracer() Option {
	return func(o *options) {
		o.splitTracer = true
	}
}

func Interceptors(opts ...Option) []interceptor.WorkerInterceptor {
	o := &options{
		tracerProvider: noop.NewTracerProvider(),
	}

	for _, opt := range opts {
		opt(o)
	}

	var tracer interceptor.WorkerInterceptor
	if o.splitTracer {
		tracer = interceptors.NewWorkerTracer(o.tracerProvider)
	} else {
		tracer = interceptors.NewTracer(o.tracerProvider)
	}

	return []interceptor.WorkerInterceptor{
		interceptors.NewActivityError(),
		tracer,
	}
}
