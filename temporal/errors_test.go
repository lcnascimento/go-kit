package temporal_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	sdktemporal "go.temporal.io/sdk/temporal"

	"github.com/lcnascimento/go-kit/errors"

	"github.com/lcnascimento/go-kit/temporal"
)

func TestFailureReason(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "nil error",
			err:  nil,
			want: "",
		},
		{
			name: "non-Temporal error",
			err:  errors.New("generic error"),
			want: "generic error",
		},
		{
			name: "ApplicationError with message",
			err:  sdktemporal.NewApplicationError("app error message", "type1"),
			want: "app error message",
		},
		{
			name: "wrapped error",
			err:  fmt.Errorf("wrapped error: %w", errors.New("wrapped error message").WithCode("TESTE")),
			want: "TESTE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := temporal.FailureReason(tt.err)
			assert.Equal(t, tt.want, got, "FailureReason() should return the same value as the error's message")
		})
	}
}
