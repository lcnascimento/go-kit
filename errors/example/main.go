package sample

import (
	"fmt"

	"github.com/lcnascimento/go-kit/errors"
)

// For more examples on how to use this package, please read the tests.
func main() {
	// basic usage
	errors.New("basic error")
	errors.New("error with kind").WithKind(errors.KindNotFound)
	errors.New("error with code").WithCode("MY_CUSTOM_CODE")
	errors.New("error with cause").WithCause(fmt.Errorf("root error"))
	errors.New("error retryable").Retryable()
	errors.New("error fully filled").
		WithKind(errors.KindInternal).
		WithCode("ERR_GET_ACCOUNT").
		WithCause(fmt.Errorf("connection refused")).
		Retryable()

	// extracting error data
	errors.Code(errors.ErrResourceNotFound)        // CodeUnknown
	errors.Kind(errors.ErrResourceNotFound)        // KindNotFound
	errors.IsRetryable(errors.ErrResourceNotFound) // false

	// go-like utility features
	wrapped := errors.Wrap(errors.ErrResourceNotFound, "could not find the requested account")
	errors.Unwrap(wrapped)                                          // ["could not find the requested account", "resource not found"]
	errors.Is(errors.ErrResourceNotFound, errors.ErrNotImplemented) // false
	errors.As(errors.ErrResourceNotFound, new(errors.CustomError))  // true
}
