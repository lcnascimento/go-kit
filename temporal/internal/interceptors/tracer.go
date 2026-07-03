package interceptors

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/workflow"
)

type spanKeyType struct{}

var spanKey = spanKeyType{}

// SpanFromWorkflowContext extracts the OTel span stored in a workflow.Context
// by the tracing interceptor. Returns nil if no span is found.
func SpanFromWorkflowContext(ctx workflow.Context) trace.Span {
	val := ctx.Value(spanKey)
	if val == nil {
		return nil
	}

	span, ok := val.(trace.Span)
	if !ok {
		return nil
	}

	return span
}

func newTracingBase(tp trace.TracerProvider) interceptor.Interceptor {
	base, err := opentelemetry.NewTracingInterceptor(opentelemetry.TracerOptions{
		Tracer:            tp.Tracer("github.com/lcnascimento/go-kit/temporal"),
		TextMapPropagator: otel.GetTextMapPropagator(),
		SpanContextKey:    spanKey,
	})
	if err != nil {
		panic(err)
	}

	return base
}

// tracer implements both ClientInterceptor and WorkerInterceptor.
// Kept for backward compatibility with services that haven't migrated
// to the split tracer pattern.
type tracer struct {
	base interceptor.Interceptor

	interceptor.ClientInterceptorBase
	interceptor.WorkerInterceptorBase
}

func NewTracer(tp trace.TracerProvider) interceptor.Interceptor {
	return &tracer{base: newTracingBase(tp)}
}

func (t *tracer) InterceptClient(next interceptor.ClientOutboundInterceptor) interceptor.ClientOutboundInterceptor {
	return t.base.InterceptClient(next)
}

func (t *tracer) InterceptActivity(ctx context.Context, next interceptor.ActivityInboundInterceptor) interceptor.ActivityInboundInterceptor {
	return t.base.InterceptActivity(ctx, next)
}

func (t *tracer) InterceptWorkflow(ctx workflow.Context, next interceptor.WorkflowInboundInterceptor) interceptor.WorkflowInboundInterceptor {
	enricher := &tracerWorkflowInterceptor{next: next}

	return t.base.InterceptWorkflow(ctx, enricher)
}

func (t *tracer) InterceptNexusOperation(
	ctx context.Context,
	next interceptor.NexusOperationInboundInterceptor,
) interceptor.NexusOperationInboundInterceptor {
	return t.base.InterceptNexusOperation(ctx, next)
}
