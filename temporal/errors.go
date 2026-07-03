package temporal

import (
	goerrors "errors"

	sdk "go.temporal.io/sdk/temporal"

	"github.com/lcnascimento/go-kit/errors"
)

var (
	ErrStartWorkflow = errors.New("failed to start workflow").
				WithCode("ERR_START_WORKFLOW").
				WithKind(errors.KindInternal).
				Retryable()

	ErrDescribeWorkflow = errors.New("failed to describe workflow").
				WithCode("ERR_DESCRIBE_WORKFLOW").
				WithKind(errors.KindInternal).
				Retryable()

	ErrMissingWorkflowExecutionInfo = errors.New("missing workflow execution info").
					WithCode("ERR_MISSING_WORKFLOW_EXECUTION_INFO").
					WithKind(errors.KindCritical)
)

// FailureReason tries to extract a concise, human-friendly reason from Temporal error wrappers.
// It prioritizes unwrapped domain errors, then ApplicationError messages, and finally falls back to the original error.
func FailureReason(err error) string {
	if err == nil {
		return ""
	}

	unwrappedErr := goerrors.Unwrap(err)

	if unwrappedErr != nil {
		customErr := errors.New("")

		if errors.As(unwrappedErr, &customErr) {
			return string(errors.Code(customErr))
		}

		return unwrappedErr.Error()
	}

	var appErr *sdk.ApplicationError
	if errors.As(err, &appErr) {
		return appErr.Message()
	}

	return err.Error()
}

func FailureCause(err error) error {
	if err == nil {
		return nil
	}

	var appErr *sdk.ApplicationError
	if errors.As(err, &appErr) {
		return appErr.Unwrap()
	}

	return err
}
