package errors

import (
	e "errors"
	"fmt"
)

// CustomError is a structure that encodes useful information about a given error.
//
// Kind: Gives semantics for the error. It is expected to be interpreted by transport layers;
// Code: Defines what the Error actually is, by an unique alias;
// Retryable: Indicates if the given error may be fixed with a retry execution.
//
// It is designed to work well within a Go Error Tree.
type CustomError struct {
	kind      KindType
	code      CodeType
	retryable bool
	errs      []error
}

// New returns a new instance of CustomError with the given message.
// It uses KindUnknown, CodeUnknown and 'false' by default for Kind, Code and Retryable attributes, respectively.
func New(msg string, args ...any) CustomError {
	err := fmt.Errorf(msg, args...)

	return CustomError{
		kind:      KindUnknown,
		code:      CodeUnknown,
		errs:      []error{err},
		retryable: false,
	}
}

// NewMissingRequiredDependency creates a new error that indicates a missing required dependency.
// It should be producing at struct constructors.
func NewMissingRequiredDependency(name string) error {
	return New("Missing required dependency: %s", name).
		WithKind(KindInvalidInput).
		WithCode("MISSING_REQUIRED_DEPENDENCY")
}

// NewValidationError creates a Validation error.
func NewValidationError(desc string) error {
	return New(desc).WithKind(KindInvalidInput).WithCode("VALIDATION_ERROR")
}

// Error returns CustomError message.
func (ce CustomError) Error() string {
	return ce.errs[0].Error()
}

// Is indicates if the current error is equal to the given target one.
func (ce CustomError) Is(target error) bool {
	var te CustomError
	if !e.As(target, &te) {
		return false
	}

	if len(ce.errs) != len(te.errs) {
		return false
	}

	for i := range ce.errs {
		if ce.errs[i] != te.errs[i] {
			return false
		}
	}

	eq1 := ce.code == te.code
	eq2 := ce.kind == te.kind
	eq3 := ce.retryable == te.retryable

	return eq1 && eq2 && eq3
}

// Unwrap unwraps all internal errors that are baselines for this error.
func (ce CustomError) Unwrap() []error {
	return ce.errs
}

// WithKind return a copy of the CustomError with the given KindType filled.
func (ce CustomError) WithKind(kind KindType) CustomError {
	ce.kind = kind

	return ce
}

// WithCode return a copy of the CustomError with the given CodeType filled.
func (ce CustomError) WithCode(code CodeType) CustomError {
	ce.code = code

	return ce
}

// WithCause return a copy of the CustomError with the given Cause attached as the
// last internal error of this CustomError.
func (ce CustomError) WithCause(cause error) CustomError {
	ce.errs = append(ce.errs, cause)

	return ce
}

// Retryable returns a copy of the CustomError tagged as retryable.
func (ce CustomError) Retryable() CustomError {
	ce.retryable = true

	return ce
}

// Kind retrieves the first non unknown Kind in err's tree.
// KindUnknown indicates that no Kind was set or no CustomError was found in the tree.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap() error
// or Unwrap() []error method. When err wraps multiple errors, Is examines err followed by a depth-first
// traversal of its children.
func Kind(err error) KindType {
	//nolint:errorlint // we don't want to use [errors.As] here intentionally.
	if ce, ok := err.(CustomError); ok && ce.kind != KindUnknown {
		return ce.kind
	}

	for _, inner := range Unwrap(err) {
		if k := Kind(inner); k != KindUnknown {
			return k
		}
	}

	return KindUnknown
}

// Code retrieves the first non unknown Code in err's tree.
// CodeUnknown indicates that no Code was set or no CustomError was found in the tree.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap() error
// or Unwrap() []error method. When err wraps multiple errors, Is examines err followed by a depth-first
// traversal of its children.
func Code(err error) CodeType {
	//nolint:errorlint // we don't want to use [errors.As] here intentionally.
	if ce, ok := err.(CustomError); ok && ce.code != CodeUnknown {
		return ce.code
	}

	for _, inner := range Unwrap(err) {
		if c := Code(inner); c != CodeUnknown {
			return c
		}
	}

	return CodeUnknown
}

// IsRetryable reports whether any error in err's tree is retryable.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap() error
// or Unwrap() []error method. When err wraps multiple errors, Is examines err followed by a depth-first
// traversal of its children.
func IsRetryable(err error) bool {
	//nolint:errorlint // we don't want to use [errors.As] here intentionally.
	if ce, ok := err.(CustomError); ok && ce.retryable {
		return true
	}

	for _, inner := range Unwrap(err) {
		if IsRetryable(inner) {
			return true
		}
	}

	return false
}

// Is reports whether any error in err's tree matches target.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap() error
// or Unwrap() []error method. When err wraps multiple errors, Is examines err followed by a depth-first
// traversal of its children.
//
// An error is considered to match a target if it is equal to that target or if it implements a method
// Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent to an existing error.
// For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See [syscall.Errno.Is] for
// an example in the standard library. An Is method should only shallowly
// compare err and the target and not call [Unwrap] on either.
func Is(err, target error) bool {
	return e.Is(err, target)
}

// As finds the first error in err's tree that matches target, and if one is found, sets target to that
// error value and returns true. Otherwise, it returns false.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap() error
// or Unwrap() []error method. When err wraps multiple errors, As examines err followed by a depth-first
// traversal of its children.
//
// An error matches target if the error's concrete value is assignable to the value pointed to by target,
// or if the error has a method As(interface{}) bool such that As(target) returns true. In the latter case,
// the As method is responsible for setting target.
//
// An error type might provide an As method so it can be treated as if it were a different error type.
//
// As panics if target is not a non-nil pointer to either a type that implements error, or to any interface type.
func As(err error, target any) bool {
	return e.As(err, target)
}

// Wrap adds more contextual information into the given error.
// The new information is handled as a new error which wraps the given error properly.
func Wrap(err error, msg string, args ...any) error {
	toWrap := fmt.Errorf(msg, args...)

	//nolint:errorlint // we don't want to use [errors.As] here intentionally.
	if ce, ok := err.(CustomError); ok {
		ce.errs = append([]error{toWrap}, ce.errs...)

		return ce
	}

	return e.Join(toWrap, err)
}

// Unwrap retrieves all the errors that forms the baseline for the given one.
// It only works in the surface of the error. In other words, it does not
// evaluates all the err's tree, just its first level.
func Unwrap(err error) []error {
	if u, ok := err.(interface{ Unwrap() error }); ok {
		if inner := u.Unwrap(); inner != nil {
			return []error{inner}
		}

		return nil
	}

	if u, ok := err.(interface{ Unwrap() []error }); ok {
		return u.Unwrap()
	}

	return nil
}
