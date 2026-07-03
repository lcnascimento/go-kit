package interceptors

import (
	"context"

	"go.temporal.io/sdk/interceptor"
)

type activity struct {
	interceptor.WorkerInterceptorBase
}

func NewActivityError() interceptor.WorkerInterceptor {
	return &activity{}
}

func (i *activity) InterceptActivity(ctx context.Context, next interceptor.ActivityInboundInterceptor) interceptor.ActivityInboundInterceptor {
	return &activityInboundInterceptor{next: next}
}

type activityInboundInterceptor struct {
	next interceptor.ActivityInboundInterceptor

	interceptor.ActivityInboundInterceptorBase
}

func (i *activityInboundInterceptor) Init(outbound interceptor.ActivityOutboundInterceptor) error {
	return i.next.Init(outbound)
}

func (i *activityInboundInterceptor) ExecuteActivity(ctx context.Context, in *interceptor.ExecuteActivityInput) (any, error) {
	rst, err := i.next.ExecuteActivity(ctx, in)
	if err != nil {
		onActivityError(ctx, err)
	}

	return rst, err
}
