package errors_test

import (
	e "errors"
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
			assert.Equal(t, errors.KindUnexpected, errors.Kind(err))
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

	t.Run("should produce an error with a RootError attached", func(t *testing.T) {
		rootErr := e.New("root error")
		err := errors.New("mocked message").WithRootError(rootErr)

		assert.Equal(t, rootErr.Error(), errors.RootError(err))
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
		assert.Equal(t, false, errors.Retryable(err))
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
		assert.Equal(t, false, errors.Retryable(err))
	})
}

func TestKind(t *testing.T) {
	tt := []struct {
		name         string
		err          error
		expectedKind errors.KindType
	}{
		{
			name:         "go native error",
			err:          e.New("new error"),
			expectedKind: errors.KindUnexpected,
		},
		{
			name:         "custom error with default kind",
			err:          errors.New("some message"),
			expectedKind: errors.KindUnexpected,
		},
		{
			name:         "custom error with non-default kind",
			err:          errors.New("some message").WithKind("some kind"),
			expectedKind: "some kind",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedKind, errors.Kind(tc.err))
		})
	}
}

func TestCode(t *testing.T) {
	tt := []struct {
		name         string
		err          error
		expectedCode errors.CodeType
	}{
		{
			name:         "go native error",
			err:          e.New("new error"),
			expectedCode: errors.CodeUnknown,
		},
		{
			name:         "custom error with default code",
			err:          errors.New("some message"),
			expectedCode: errors.CodeUnknown,
		},
		{
			name:         "custom error with non-default code",
			err:          errors.New("some message").WithCode("some code"),
			expectedCode: "some code",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedCode, errors.Code(tc.err))
		})
	}
}

func TestRootError(t *testing.T) {
	tt := []struct {
		name               string
		err                error
		expectedErrMsg     string
		expectedRootErrMsg string
	}{
		{
			name:               "go native error",
			err:                errors.New("go native error"),
			expectedErrMsg:     "go native error",
			expectedRootErrMsg: "go native error",
		},
		{
			name:               "custom error with attached root error",
			err:                errors.New("custom error").WithRootError(e.New("root error")),
			expectedErrMsg:     "custom error",
			expectedRootErrMsg: "root error",
		},
		{
			name:               "custom error without attached root error",
			err:                errors.New("custom error"),
			expectedErrMsg:     "custom error",
			expectedRootErrMsg: "custom error",
		},
		{
			name:               "chain of root errors",
			err:                errors.New("head error").WithRootError(errors.New("middle error").WithRootError(errors.New("tail error"))),
			expectedErrMsg:     "head error",
			expectedRootErrMsg: "tail error",
		},
		{
			name:               "nil error",
			err:                nil,
			expectedErrMsg:     "",
			expectedRootErrMsg: "",
		},
		{
			name:               "error without message and root error",
			err:                errors.New(""),
			expectedErrMsg:     "",
			expectedRootErrMsg: "",
		},
		{
			name:               "error without message but with root error",
			err:                errors.New("").WithRootError(e.New("root error")),
			expectedErrMsg:     "root error",
			expectedRootErrMsg: "root error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var err string
			if tc.err != nil {
				err = tc.err.Error()
			}

			rootErr := errors.RootError(tc.err)

			t.Run("should produce right error message", func(t *testing.T) {
				assert.Equal(t, tc.expectedErrMsg, err)
			})

			t.Run("should produce right root error message", func(t *testing.T) {
				assert.Equal(t, tc.expectedRootErrMsg, rootErr)
			})
		})
	}
}

func TestRetryable(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "go native error",
			err:      e.New("new error"),
			expected: false,
		},
		{
			name:     "custom error with default retry mode",
			err:      errors.New("some message"),
			expected: false,
		},
		{
			name:     "custom error with non-default retry mode",
			err:      errors.New("some message").Retryable(true),
			expected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, errors.Retryable(tc.err))
		})
	}
}

func TestIs(t *testing.T) {
	t.Run("should see different errors", func(t *testing.T) {
		assert.False(t, errors.Is(errors.ErrNotImplemented, errors.ErrResourceNotFound))
	})

	t.Run("should see equal errors", func(t *testing.T) {
		assert.True(t, errors.Is(errors.ErrNotImplemented, errors.ErrNotImplemented))
	})
}
