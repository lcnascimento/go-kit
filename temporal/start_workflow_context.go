package temporal

import (
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

const workflowExecutionTimeout = 4 * 24 * time.Hour

type startWorkflowContextOptions struct {
	workflowIDReusePolicy                    enums.WorkflowIdReusePolicy
	workflowIDConflictPolicy                 enums.WorkflowIdConflictPolicy
	workflowExecutionTimeout                 time.Duration
	workflowExecutionErrorWhenAlreadyStarted bool
	searchAttributes                         map[string]string
}

type StartWorkflowContextOption func(*startWorkflowContextOptions)

func WithWorkflowIDReusePolicy(policy enums.WorkflowIdReusePolicy) StartWorkflowContextOption {
	return func(o *startWorkflowContextOptions) {
		o.workflowIDReusePolicy = policy
	}
}

func WithWorkflowIDConflictPolicy(policy enums.WorkflowIdConflictPolicy) StartWorkflowContextOption {
	return func(o *startWorkflowContextOptions) {
		o.workflowIDConflictPolicy = policy
	}
}

func WithWorkflowExecutionTimeout(timeout time.Duration) StartWorkflowContextOption {
	return func(o *startWorkflowContextOptions) {
		o.workflowExecutionTimeout = timeout
	}
}

func WithWorkflowExecutionErrorWhenAlreadyStarted() StartWorkflowContextOption {
	return func(o *startWorkflowContextOptions) {
		o.workflowExecutionErrorWhenAlreadyStarted = true
	}
}

func WithSearchAttributeOrganization(org string) StartWorkflowContextOption {
	return func(o *startWorkflowContextOptions) {
		o.searchAttributes["organization"] = org
	}
}

func WithSearchAttributeUserID(userID string) StartWorkflowContextOption {
	return func(o *startWorkflowContextOptions) {
		o.searchAttributes["user_id"] = userID
	}
}

func WithSearchAttributeAccountID(accountID uuid.UUID) StartWorkflowContextOption {
	return func(o *startWorkflowContextOptions) {
		o.searchAttributes["account_id"] = accountID.String()
	}
}

func StartWorkflowOptions(wID, queue string, opts ...StartWorkflowContextOption) client.StartWorkflowOptions {
	options := &startWorkflowContextOptions{
		workflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
		workflowIDConflictPolicy:                 enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
		workflowExecutionTimeout:                 workflowExecutionTimeout,
		workflowExecutionErrorWhenAlreadyStarted: false,
		searchAttributes:                         make(map[string]string, 0),
	}

	for _, opt := range opts {
		opt(options)
	}

	attrs := make([]temporal.SearchAttributeUpdate, 0)
	for k, v := range options.searchAttributes {
		attrs = append(attrs, temporal.NewSearchAttributeKeyString(k).ValueSet(v))
	}

	searchAttributes := temporal.NewSearchAttributes(attrs...)

	return client.StartWorkflowOptions{
		ID:                                       wID,
		TaskQueue:                                queue,
		WorkflowIDReusePolicy:                    options.workflowIDReusePolicy,
		WorkflowExecutionTimeout:                 options.workflowExecutionTimeout,
		WorkflowIDConflictPolicy:                 options.workflowIDConflictPolicy,
		WorkflowExecutionErrorWhenAlreadyStarted: options.workflowExecutionErrorWhenAlreadyStarted,
		TypedSearchAttributes:                    searchAttributes,
	}
}
