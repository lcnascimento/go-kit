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

func TestSeverity(t *testing.T) {
	t.Run("go native error", func(t *testing.T) {
		err := e.New("new error")
		assert.Equal(t, errors.SeverityError, errors.Severity(err))
	})

	t.Run("custom error without kind", func(t *testing.T) {
		err := errors.New("some message")
		assert.Equal(t, errors.SeverityError, errors.Severity(err))
	})

	t.Run("custom error with kind InvalidInput", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindInvalidInput)
		assert.Equal(t, errors.SeverityWarn, errors.Severity(err))
	})

	t.Run("custom error with kind NotFound", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindNotFound)
		assert.Equal(t, errors.SeverityWarn, errors.Severity(err))
	})

	t.Run("custom error with kind Unauthorized", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindUnauthorized)
		assert.Equal(t, errors.SeverityError, errors.Severity(err))
	})

	t.Run("custom error with kind Unauthenticated", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindUnauthenticated)
		assert.Equal(t, errors.SeverityError, errors.Severity(err))
	})

	t.Run("custom error with kind ResourceExhausted", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindResourceExhausted)
		assert.Equal(t, errors.SeverityError, errors.Severity(err))
	})

	t.Run("custom error with kind Conflict", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindConflict)
		assert.Equal(t, errors.SeverityError, errors.Severity(err))
	})

	t.Run("custom error with kind Internal", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindInternal)
		assert.Equal(t, errors.SeverityError, errors.Severity(err))
	})

	t.Run("custom error with kind Critical", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindCritical)
		assert.Equal(t, errors.SeverityCritical, errors.Severity(err))
	})

	t.Run("custom error with kind Fatal", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindFatal)
		assert.Equal(t, errors.SeverityFatal, errors.Severity(err))
	})

	t.Run("wrap custom error", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindNotFound)
		wrapped := fmt.Errorf("wrapped: %w", err)
		assert.Equal(t, errors.SeverityWarn, errors.Severity(wrapped))
	})

	t.Run("join native error with custom error", func(t *testing.T) {
		native := e.New("native error")
		custom := errors.New("some message").WithKind(errors.KindCritical)
		wrapped := e.Join(native, custom)

		assert.Equal(t, errors.SeverityCritical, errors.Severity(wrapped))
	})

	t.Run("join two custom errors", func(t *testing.T) {
		custom1 := errors.New("custom error 1")
		custom2 := errors.New("custom error 2").WithKind(errors.KindCritical)
		wrapped := e.Join(custom1, custom2)

		assert.Equal(t, errors.SeverityCritical, errors.Severity(wrapped))
	})

	t.Run("custom error caused by other custom error", func(t *testing.T) {
		custom1 := errors.New("custom error 1").WithKind(errors.KindCritical)
		custom2 := errors.New("custom error 2").WithCause(custom1)

		assert.Equal(t, errors.SeverityCritical, errors.Severity(custom2))
	})

	t.Run("severity WARN name error", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindNotFound)
		assert.Equal(t, "WARN", errors.Severity(err).String())
	})

	t.Run("severity ERROR name error", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindInternal)
		assert.Equal(t, "ERROR", errors.Severity(err).String())
	})

	t.Run("severity CRITICAL name error", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindCritical)
		assert.Equal(t, "CRITICAL", errors.Severity(err).String())
	})

	t.Run("severity FATAL name error", func(t *testing.T) {
		err := errors.New("some message").WithKind(errors.KindFatal)
		assert.Equal(t, "FATAL", errors.Severity(err).String())
	})
}

func TestAttributes(t *testing.T) {
	t.Run("go native error", func(t *testing.T) {
		err := e.New("new error")
		assert.Equal(t, errors.AttributeSet{}, errors.Attributes(err))
	})

	t.Run("custom error with no attributes", func(t *testing.T) {
		err := errors.New("some message")
		assert.Equal(t, errors.AttributeSet{}, errors.Attributes(err))
	})

	t.Run("custom error with basic attributes", func(t *testing.T) {
		err := errors.New("some message").WithAttribute("key", "value")
		assert.Equal(t, errors.AttributeSet{"key": "value"}, errors.Attributes(err))
	})

	t.Run("custom error with nested attributes", func(t *testing.T) {
		err1 := errors.New("nested").WithAttribute("key", "value")
		err2 := errors.New("nested").WithAttribute("key2", "value2").WithCause(err1)
		wrapped := errors.New("error").WithCause(err2)

		assert.Equal(t, errors.AttributeSet{"key": "value", "key2": "value2"}, errors.Attributes(wrapped))
	})

	t.Run("custom error with merged attributes", func(t *testing.T) {
		err := errors.New("some message").WithAttribute("key", "value")
		err = err.WithAttribute("key2", "value2")

		assert.Equal(t, errors.AttributeSet{"key": "value", "key2": "value2"}, errors.Attributes(err))
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

func TestReasons(t *testing.T) {
	t.Run("native error", func(t *testing.T) {
		err := e.New("some message")
		assert.Empty(t, errors.Reasons(err))
	})

	t.Run("custom error with no reasons", func(t *testing.T) {
		err := errors.New("some message")
		assert.Empty(t, errors.Reasons(err))
	})

	t.Run("native error with wrapped errors", func(t *testing.T) {
		base := e.New("base")
		wrap1 := e.New("wrap1")
		wrap2 := e.New("wrap2")
		err := e.Join(base, wrap1)
		err = e.Join(err, wrap2)

		assert.Equal(t, []string{"wrap1", "wrap2"}, errors.Reasons(err))
	})

	t.Run("custom error with native reason", func(t *testing.T) {
		reason := e.New("reason")
		err := errors.New("some message").WithCause(reason)
		assert.Equal(t, []string{"reason"}, errors.Reasons(err))
	})

	t.Run("custom error with custom and native reasons", func(t *testing.T) {
		reason1 := e.New("reason 1")
		reason2 := errors.New("reason 2")
		err := errors.New("some message").WithCause(reason1).WithCause(reason2)
		assert.Equal(t, []string{"reason 1", "reason 2"}, errors.Reasons(err))
	})

	t.Run("custom error with complex reason", func(t *testing.T) {
		reason2 := e.New("reason 2")
		reason1 := errors.New("reason 1").WithCause(reason2)
		err := errors.New("some message").WithCause(reason1)
		assert.Equal(t, []string{"reason 1", "reason 2"}, errors.Reasons(err))
	})

	t.Run("native error with complex reason", func(t *testing.T) {
		reason2 := e.New("reason 2")
		reason1 := errors.New("reason 1").WithCause(reason2)
		err := e.Join(e.New("some message"), reason1)
		assert.Equal(t, []string{"reason 1", "reason 2"}, errors.Reasons(err))
	})
}

func TestSafeReasons(t *testing.T) {
	t.Run("native error", func(t *testing.T) {
		err := e.New("some message")
		assert.Empty(t, errors.SafeReasons(err))
	})

	t.Run("custom error with no reasons", func(t *testing.T) {
		err := errors.New("some message")
		assert.Empty(t, errors.SafeReasons(err))
	})

	t.Run("native error with wrapped errors", func(t *testing.T) {
		base := e.New("base")
		wrap1 := e.New("wrap1")
		wrap2 := e.New("wrap2")
		err := e.Join(base, wrap1)
		err = e.Join(err, wrap2)

		assert.Equal(t, []string{}, errors.SafeReasons(err))
	})

	t.Run("custom error with native reason", func(t *testing.T) {
		reason := e.New("reason")
		err := errors.New("some message").WithCause(reason)
		assert.Equal(t, []string{}, errors.SafeReasons(err))
	})

	t.Run("custom error with custom and native reasons", func(t *testing.T) {
		reason1 := e.New("reason 1")
		reason2 := errors.New("reason 2")
		err := errors.New("some message").WithCause(reason1).WithCause(reason2)
		assert.Equal(t, []string{"reason 2"}, errors.SafeReasons(err))
	})

	t.Run("custom error with complex reason", func(t *testing.T) {
		reason2 := e.New("reason 2")
		reason1 := errors.New("reason 1").WithCause(reason2)
		err := errors.New("some message").WithCause(reason1)
		assert.Equal(t, []string{"reason 1"}, errors.SafeReasons(err))
	})

	t.Run("native error with complex reason", func(t *testing.T) {
		reason2 := e.New("reason 2")
		reason1 := errors.New("reason 1").WithCause(reason2)
		err := e.Join(e.New("some message"), reason1)
		assert.Equal(t, []string{"reason 1"}, errors.SafeReasons(err))
	})
}

func TestIs(t *testing.T) {
	t.Run("go native and custom error", func(t *testing.T) {
		assert.False(t, errors.Is(errors.ErrResourceNotFound, e.New("go native")))
	})

	t.Run("different custom errors", func(t *testing.T) {
		assert.False(t, errors.Is(errors.ErrNotImplemented, errors.ErrResourceNotFound))
	})

	t.Run("custom errors with equal internal errors", func(t *testing.T) {
		native := fmt.Errorf("go native error")
		custom1 := errors.New("custom error").WithCause(native)
		custom2 := errors.New("custom error").WithCause(native)

		assert.True(t, errors.Is(custom1, custom2))
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

	t.Run("join native with custom error", func(t *testing.T) {
		cause := e.New("cause error")
		native := e.New("native error")
		custom := errors.New("custom error").WithKind(errors.KindNotFound).WithCause(cause)
		wrapped := e.Join(native, custom)

		assert.True(t, e.Is(wrapped, native))
		assert.True(t, e.Is(wrapped, custom))
		assert.True(t, e.Is(wrapped, cause))
	})

	t.Run("with cause", func(t *testing.T) {
		base := errors.New("base error")
		cause := e.New("cause error")
		final := base.WithCause(cause)

		assert.True(t, errors.Is(final, base), "should be equivalent to base")
		assert.True(t, errors.Is(final, cause), "should be equivalent to cause")
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
