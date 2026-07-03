package temporal

import (
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const childWorkflowExecutionTimeout = 4 * 24 * time.Hour

type childWorkflowContextOptions struct {
	taskQueue                string
	parentClosePolicy        enums.ParentClosePolicy
	workflowIDReusePolicy    enums.WorkflowIdReusePolicy
	workflowExecutionTimeout time.Duration
	searchAttributes         map[string]string
}

type ChildWorkflowContextOption func(*childWorkflowContextOptions)

func WithChildWorkflowTaskQueue(queue string) ChildWorkflowContextOption {
	return func(o *childWorkflowContextOptions) {
		o.taskQueue = queue
	}
}

func WithChildWorkflowParentClosePolicy(policy enums.ParentClosePolicy) ChildWorkflowContextOption {
	return func(o *childWorkflowContextOptions) {
		o.parentClosePolicy = policy
	}
}

func WithChildWorkflowWorkflowIDReusePolicy(policy enums.WorkflowIdReusePolicy) ChildWorkflowContextOption {
	return func(o *childWorkflowContextOptions) {
		o.workflowIDReusePolicy = policy
	}
}

func WithChildWorkflowWorkflowExecutionTimeout(timeout time.Duration) ChildWorkflowContextOption {
	return func(o *childWorkflowContextOptions) {
		o.workflowExecutionTimeout = timeout
	}
}

func WithChildWorkflowSearchAttributeOrganization(org string) ChildWorkflowContextOption {
	return func(o *childWorkflowContextOptions) {
		o.searchAttributes["organization"] = org
	}
}

func WithChildWorkflowSearchAttributeUserID(userID string) ChildWorkflowContextOption {
	return func(o *childWorkflowContextOptions) {
		o.searchAttributes["user_id"] = userID
	}
}

func WithChildWorkflowSearchAttributeAccountID(accountID uuid.UUID) ChildWorkflowContextOption {
	return func(o *childWorkflowContextOptions) {
		o.searchAttributes["account_id"] = accountID.String()
	}
}

func WithChildWorkflowIDReusePolicy(policy enums.WorkflowIdReusePolicy) ChildWorkflowContextOption {
	return func(o *childWorkflowContextOptions) {
		o.workflowIDReusePolicy = policy
	}
}

func ChildWorkflowContext(ctx workflow.Context, wID string, opts ...ChildWorkflowContextOption) workflow.Context {
	options := &childWorkflowContextOptions{
		parentClosePolicy:        enums.PARENT_CLOSE_POLICY_TERMINATE,
		workflowIDReusePolicy:    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
		workflowExecutionTimeout: childWorkflowExecutionTimeout,
		searchAttributes:         make(map[string]string, 0),
	}

	for _, opt := range opts {
		opt(options)
	}

	attrs := make([]temporal.SearchAttributeUpdate, 0)
	for k, v := range options.searchAttributes {
		attrs = append(attrs, temporal.NewSearchAttributeKeyString(k).ValueSet(v))
	}

	searchAttributes := temporal.NewSearchAttributes(attrs...)

	return workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID:               wID,
		TaskQueue:                options.taskQueue,
		WorkflowIDReusePolicy:    options.workflowIDReusePolicy,
		WorkflowExecutionTimeout: options.workflowExecutionTimeout,
		ParentClosePolicy:        options.parentClosePolicy,
		TypedSearchAttributes:    searchAttributes,
	})
}
