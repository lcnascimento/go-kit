package temporal

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	activityStartToCloseTimeout        = time.Minute
	activityScheduleToCloseTimeout     = 24 * time.Hour
	activityInitialInterval            = time.Second
	activityMaximumInterval            = 5 * time.Minute
	activityHeartbeatTimeout           = 2 * time.Minute
	activityPollingStartToCloseTimeout = 24 * time.Hour
)

type activityContextOptions struct {
	startToCloseTimeout    time.Duration
	scheduleToCloseTimeout time.Duration
	initialInterval        time.Duration
	maximumInterval        time.Duration
	heartbeatTimeout       time.Duration
	nonRetryableErrorTypes []string
}

type ActivityContextOption func(*activityContextOptions)

func WithActivityStartToCloseTimeout(timeout time.Duration) ActivityContextOption {
	return func(o *activityContextOptions) {
		o.startToCloseTimeout = timeout
	}
}

func WithActivityScheduleToCloseTimeout(timeout time.Duration) ActivityContextOption {
	return func(o *activityContextOptions) {
		o.scheduleToCloseTimeout = timeout
	}
}

func WithActivityInitialInterval(interval time.Duration) ActivityContextOption {
	return func(o *activityContextOptions) {
		o.initialInterval = interval
	}
}

func WithActivityMaximumInterval(interval time.Duration) ActivityContextOption {
	return func(o *activityContextOptions) {
		o.maximumInterval = interval
	}
}

func WithActivityHeartbeatTimeout(timeout time.Duration) ActivityContextOption {
	return func(o *activityContextOptions) {
		o.heartbeatTimeout = timeout
	}
}

func WithPolling() ActivityContextOption {
	return func(o *activityContextOptions) {
		o.startToCloseTimeout = activityPollingStartToCloseTimeout
	}
}

func WithActivityNonRetryableErrorTypes(types []string) ActivityContextOption {
	return func(o *activityContextOptions) {
		o.nonRetryableErrorTypes = types
	}
}

func ActivityContext(ctx workflow.Context, queue string, opts ...ActivityContextOption) workflow.Context {
	options := &activityContextOptions{
		startToCloseTimeout:    activityStartToCloseTimeout,
		scheduleToCloseTimeout: activityScheduleToCloseTimeout,
		initialInterval:        activityInitialInterval,
		maximumInterval:        activityMaximumInterval,
		heartbeatTimeout:       activityHeartbeatTimeout,
	}

	for _, opt := range opts {
		opt(options)
	}

	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              queue,
		StartToCloseTimeout:    options.startToCloseTimeout,
		ScheduleToCloseTimeout: options.scheduleToCloseTimeout,
		HeartbeatTimeout:       options.heartbeatTimeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        options.initialInterval,
			MaximumInterval:        options.maximumInterval,
			NonRetryableErrorTypes: options.nonRetryableErrorTypes,
		},
	})
}
