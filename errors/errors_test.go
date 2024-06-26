package errors_test

import (
	e "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lcnascimento/go-kit/errors"
)

func TestNew(t *testing.T) {
	t.Run("should produce an error with the given message", func(t *testing.T) {
		err := errors.New("mocked message")

		assert.Equal(t, "mocked message", err.Error())
	})

	t.Run("should produce an error with an dynamic message", func(t *testing.T) {
		err := errors.New("mocked message with dynamic value %d", 1000)

		assert.Equal(t, "mocked message with dynamic value 1000", err.Error())
	})

	t.Run("should produce an error with the default values", func(t *testing.T) {
		err := errors.New("mocked message")

		t.Run("should have Kind unexpected", func(t *testing.T) {
			assert.Equal(t, errors.KindUnknown, errors.Kind(err))
		})

		t.Run("should have Code unknown", func(t *testing.T) {
			assert.Equal(t, errors.CodeUnknown, errors.Code(err))
		})
	})

	t.Run("should produce an error with a non-default Kind filled", func(t *testing.T) {
		err := errors.New("mocked message").WithKind(errors.KindInternal)

		assert.Equal(t, errors.KindInternal, errors.Kind(err))
	})

	t.Run("should override the Kind if WithKind is called more than once", func(t *testing.T) {
		err := errors.New("mocked message").WithKind(errors.KindInternal).WithKind(errors.KindNotFound)

		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("should produce an error with a non-default Code filled", func(t *testing.T) {
		err := errors.New("mocked message").WithCode("MOCKED_CODE")

		assert.Equal(t, errors.CodeType("MOCKED_CODE"), errors.Code(err))
	})

	t.Run("should override the Code if WithCode is called more than once", func(t *testing.T) {
		err := errors.New("mocked message").WithCode("MOCKED_CODE").WithCode("MOCKED_CODE_2")

		assert.Equal(t, errors.CodeType("MOCKED_CODE_2"), errors.Code(err))
	})
}

func TestNewMissingRequiredDependency(t *testing.T) {
	err := errors.NewMissingRequiredDependency("SomeDependency")

	t.Run("should produce the right error message", func(t *testing.T) {
		assert.Equal(t, "Missing required dependency: SomeDependency", err.Error())
	})

	t.Run("should produce the right error code", func(t *testing.T) {
		assert.Equal(t, errors.CodeType("MISSING_REQUIRED_DEPENDENCY"), errors.Code(err))
	})

	t.Run("should produce the right error kind", func(t *testing.T) {
		assert.Equal(t, errors.KindInvalidInput, errors.Kind(err))
	})

	t.Run("should not be retryable", func(t *testing.T) {
		assert.False(t, errors.IsRetryable(err))
	})
}

func TestNewValidationError(t *testing.T) {
	err := errors.NewValidationError("missing required field")

	t.Run("should produce the right error message", func(t *testing.T) {
		assert.Equal(t, "missing required field", err.Error())
	})

	t.Run("should produce the right error code", func(t *testing.T) {
		assert.Equal(t, errors.CodeType("VALIDATION_ERROR"), errors.Code(err))
	})

	t.Run("should produce the right error kind", func(t *testing.T) {
		assert.Equal(t, errors.KindInvalidInput, errors.Kind(err))
	})

	t.Run("should not be retryable", func(t *testing.T) {
		assert.False(t, errors.IsRetryable(err))
	})
}

func TestKind(t *testing.T) {
	t.Run("go native error", func(t *testing.T) {
		err := e.New("new error")
		assert.Equal(t, errors.KindUnknown, errors.Kind(err))
	})

	t.Run("custom error with default kind", func(t *testing.T) {
		err := errors.New("some message")
		assert.Equal(t, errors.KindUnknown, errors.Kind(err))
	})

	t.Run("custom error with non default kind", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindNotFound)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("wrap custom error", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindNotFound)
		wrapped := fmt.Errorf("wrapped: %w", err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(wrapped))
	})

	t.Run("join native error with custom error", func(t *testing.T) {
		native := e.New("native error")
		custom := errors.New("custom error").WithKind(errors.KindNotFound)
		wrapped := e.Join(native, custom)

		assert.Equal(t, errors.KindNotFound, errors.Kind(wrapped))
	})

	t.Run("join two custom errors", func(t *testing.T) {
		custom1 := errors.New("custom error 1")
		custom2 := errors.New("custom error 2").WithKind(errors.KindNotFound)
		wrapped := e.Join(custom1, custom2)

		assert.Equal(t, errors.KindNotFound, errors.Kind(wrapped))
	})

	t.Run("custom error caused by other custom error", func(t *testing.T) {
		custom1 := errors.New("custom error 1").WithKind(errors.KindNotFound)
		custom2 := errors.New("custom error 2").WithCause(custom1)

		assert.Equal(t, errors.KindNotFound, errors.Kind(custom2))
	})
}

func TestCode(t *testing.T) {
	t.Run("go native error", func(t *testing.T) {
		err := e.New("new error")
		assert.Equal(t, errors.CodeUnknown, errors.Code(err))
	})

	t.Run("custom error with default code", func(t *testing.T) {
		err := errors.New("some message")
		assert.Equal(t, errors.CodeUnknown, errors.Code(err))
	})

	t.Run("custom error with non default code", func(t *testing.T) {
		err := errors.New("some message").WithCode("SOME_CODE")
		assert.Equal(t, errors.CodeType("SOME_CODE"), errors.Code(err))
	})

	t.Run("wrap custom error", func(t *testing.T) {
		err := errors.New("some message").WithCode("SOME_CODE")
		wrapped := fmt.Errorf("wrapped: %w", err)
		assert.Equal(t, errors.CodeType("SOME_CODE"), errors.Code(wrapped))
	})

	t.Run("join native error with custom error", func(t *testing.T) {
		native := e.New("native error")
		custom := errors.New("some message").WithCode("SOME_CODE")
		wrapped := e.Join(native, custom)

		assert.Equal(t, errors.CodeType("SOME_CODE"), errors.Code(wrapped))
	})

	t.Run("join two custom errors", func(t *testing.T) {
		custom1 := errors.New("custom error 1")
		custom2 := errors.New("custom error 2").WithCode("SOME_CODE")
		wrapped := e.Join(custom1, custom2)

		assert.Equal(t, errors.CodeType("SOME_CODE"), errors.Code(wrapped))
	})

	t.Run("custom error caused by other custom error", func(t *testing.T) {
		custom1 := errors.New("custom error 1").WithCode("SOME_CODE")
		custom2 := errors.New("custom error 2").WithCause(custom1)

		assert.Equal(t, errors.CodeType("SOME_CODE"), errors.Code(custom2))
	})
}

func TestIsRetryable(t *testing.T) {
	t.Run("go native error", func(t *testing.T) {
		err := e.New("new error")
		assert.False(t, errors.IsRetryable(err))
	})

	t.Run("custom error with default retry mode", func(t *testing.T) {
		err := errors.New("some message")
		assert.False(t, errors.IsRetryable(err))
	})

	t.Run("custom error with retry mode", func(t *testing.T) {
		err := errors.New("some message").Retryable()
		assert.True(t, errors.IsRetryable(err))
	})

	t.Run("wrap custom retryable error", func(t *testing.T) {
		err := errors.New("some message").Retryable()
		wrapped := fmt.Errorf("wrapped: %w", err)
		assert.True(t, errors.IsRetryable(wrapped))
	})

	t.Run("join native error with custom retryable error", func(t *testing.T) {
		native := e.New("native error")
		custom := errors.New("some message").Retryable()
		wrapped := e.Join(native, custom)

		assert.True(t, errors.IsRetryable(wrapped))
	})

	t.Run("join two custom errors", func(t *testing.T) {
		custom1 := errors.New("custom error 1")
		custom2 := errors.New("custom error 2").Retryable()
		wrapped := e.Join(custom1, custom2)

		assert.True(t, errors.IsRetryable(wrapped))
	})

	t.Run("custom error caused by other custom retryable error", func(t *testing.T) {
		custom1 := errors.New("custom error 1").Retryable()
		custom2 := errors.New("custom error 2").WithCause(custom1)

		assert.True(t, errors.IsRetryable(custom2))
	})
}

func TestIs(t *testing.T) {
	t.Run("go native and custom error", func(t *testing.T) {
		assert.False(t, errors.Is(errors.ErrResourceNotFound, e.New("go native")))
	})

	t.Run("different custom errors", func(t *testing.T) {
		assert.False(t, errors.Is(errors.ErrNotImplemented, errors.ErrResourceNotFound))
	})

	t.Run("custom errors with different internal errors", func(t *testing.T) {
		custom1 := errors.New("custom error")
		custom2 := errors.New("custom error").WithCause(fmt.Errorf("go native error"))

		assert.False(t, errors.Is(custom1, custom2))
	})

	t.Run("equal custom errors", func(t *testing.T) {
		assert.True(t, errors.Is(errors.ErrNotImplemented, errors.ErrNotImplemented))
	})

	t.Run("different native errors", func(t *testing.T) {
		err1 := fmt.Errorf("fake error 1")
		err2 := fmt.Errorf("fake error 2")

		assert.False(t, errors.Is(err1, err2))
	})

	t.Run("equal native errors", func(t *testing.T) {
		err := fmt.Errorf("fake error")

		assert.True(t, errors.Is(err, err))
	})

	t.Run("custom error caused by native error", func(t *testing.T) {
		err := fmt.Errorf("fake error")
		custom := errors.New("some error").WithCause(err)

		assert.True(t, errors.Is(custom, err))
	})
}

func TestAs(t *testing.T) {
	t.Run("custom error", func(t *testing.T) {
		var ce errors.CustomError

		assert.True(t, errors.As(errors.ErrNotImplemented, &ce))
		assert.Equal(t, errors.ErrNotImplemented, ce)
	})
}

func TestWrap(t *testing.T) {
	t.Run("go native error", func(t *testing.T) {
		native := e.New("native error")
		wrapped := errors.Wrap(native, "wrapped error")

		expectedErrs := []error{
			fmt.Errorf("wrapped error"),
			native,
		}

		assert.Equal(t, "wrapped error\nnative error", wrapped.Error())
		assert.Equal(t, expectedErrs, errors.Unwrap(wrapped))
	})

	t.Run("custom error", func(t *testing.T) {
		custom := errors.New("custom error")
		wrapped := errors.Wrap(custom, "wrapped error")

		expectedErrs := []error{
			fmt.Errorf("wrapped error"),
			fmt.Errorf("custom error"),
		}

		assert.Equal(t, "wrapped error", wrapped.Error())
		assert.Equal(t, expectedErrs, errors.Unwrap(wrapped))
	})

	t.Run("custom error with cause", func(t *testing.T) {
		native := e.New("native error")
		custom := errors.New("custom error").WithCause(native)
		wrapped := errors.Wrap(custom, "wrapped error")

		expectedErrs := []error{
			fmt.Errorf("wrapped error"),
			fmt.Errorf("custom error"),
			native,
		}

		assert.Equal(t, "wrapped error", wrapped.Error())
		assert.Equal(t, expectedErrs, errors.Unwrap(wrapped))
	})
}
