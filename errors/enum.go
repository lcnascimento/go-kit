package errors

import "github.com/lcnascimento/go-kit/propagation"

const (
	// ContextKeyRootError defines the name of the RootError attribute attached into logs.
	ContextKeyRootError propagation.ContextKey = "error.root"

	// ContextKeyErrorKind defines the name of the ErrorKind attribute attached into logs.
	ContextKeyErrorKind propagation.ContextKey = "error.kind"

	// ContextKeyErrorCode defines the name of the ErrorCode attribute attached into logs.
	ContextKeyErrorCode propagation.ContextKey = "error.code"

	// ContextKeyErrorRetryable defines the name of the ErrorRetryable attribute attached into logs.
	ContextKeyErrorRetryable propagation.ContextKey = "error.retryable"
)
