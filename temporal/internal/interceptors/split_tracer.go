package interceptors

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/workflow"

	temporalBaggage "github.com/lcnascimento/go-kit/temporal/baggage"
)

// clientTracer implements only interceptor.ClientInterceptor.
// It must NOT implement WorkerInterceptor to prevent the Temporal SDK
// from promoting it into the worker interceptor chain.
type clientTracer struct {
	base interceptor.Interceptor

	interceptor.ClientInterceptorBase
}

func NewClientTracer(tp trace.TracerProvider) interceptor.ClientInterceptor {
	return &clientTracer{base: newTracingBase(tp)}
}

func (t *clientTracer) InterceptClient(next interceptor.ClientOutboundInterceptor) interceptor.ClientOutboundInterceptor {
	return t.base.InterceptClient(next)
}

// workerTracer implements only interceptor.WorkerInterceptor.
type workerTracer struct {
	base interceptor.Interceptor

	interceptor.WorkerInterceptorBase
}

func NewWorkerTracer(tp trace.TracerProvider) interceptor.WorkerInterceptor {
	return &workerTracer{base: newTracingBase(tp)}
}

func (t *workerTracer) InterceptActivity(ctx context.Context, next interceptor.ActivityInboundInterceptor) interceptor.ActivityInboundInterceptor {
	return t.base.InterceptActivity(ctx, next)
}

func (t *workerTracer) InterceptWorkflow(ctx workflow.Context, next interceptor.WorkflowInboundInterceptor) interceptor.WorkflowInboundInterceptor {
	enricher := &tracerWorkflowInterceptor{next: next}

	return t.base.InterceptWorkflow(ctx, enricher)
}

func (t *workerTracer) InterceptNexusOperation(
	ctx context.Context,
	next interceptor.NexusOperationInboundInterceptor,
) interceptor.NexusOperationInboundInterceptor {
	return t.base.InterceptNexusOperation(ctx, next)
}

type tracerWorkflowInterceptor struct {
	next interceptor.WorkflowInboundInterceptor

	interceptor.WorkflowInboundInterceptorBase
}

func (t *tracerWorkflowInterceptor) Init(outbound interceptor.WorkflowOutboundInterceptor) error {
	return t.next.Init(outbound)
}

func (t *tracerWorkflowInterceptor) ExecuteWorkflow(ctx workflow.Context, in *interceptor.ExecuteWorkflowInput) (any, error) {
	enrichSpanFromBaggage(ctx)

	return t.next.ExecuteWorkflow(ctx, in)
}

func (t *tracerWorkflowInterceptor) HandleSignal(ctx workflow.Context, in *interceptor.HandleSignalInput) error {
	enrichSpanFromBaggage(ctx)

	return t.next.HandleSignal(ctx, in)
}

func (t *tracerWorkflowInterceptor) HandleQuery(ctx workflow.Context, in *interceptor.HandleQueryInput) (any, error) {
	enrichSpanFromBaggage(ctx)

	return t.next.HandleQuery(ctx, in)
}

func (t *tracerWorkflowInterceptor) ExecuteUpdate(ctx workflow.Context, in *interceptor.UpdateInput) (any, error) {
	enrichSpanFromBaggage(ctx)

	return t.next.ExecuteUpdate(ctx, in)
}

func (t *tracerWorkflowInterceptor) ValidateUpdate(ctx workflow.Context, in *interceptor.UpdateInput) error {
	enrichSpanFromBaggage(ctx)

	return t.next.ValidateUpdate(ctx, in)
}

func enrichSpanFromBaggage(ctx workflow.Context) {
	spanVal := ctx.Value(spanKey)
	if spanVal == nil {
		return
	}

	span, ok := spanVal.(trace.Span)
	if !ok {
		return
	}

	bag := temporalBaggage.FromContext(ctx)

	members := bag.Members()
	if len(members) == 0 {
		return
	}

	attrs := make([]attribute.KeyValue, 0, len(members))
	for _, m := range members {
		attrs = append(attrs, attribute.String("bag."+m.Key(), m.Value()))
	}

	span.SetAttributes(attrs...)
}
