package interceptor_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/grpc/grpcserver/interceptor"
)

func TestErrorHandler(t *testing.T) {
	tt := []struct {
		desc    string
		err     errors.CustomError
		message string
		status  codes.Code
		details map[string]string
	}{
		{
			desc:    "error with code",
			err:     errors.New("test error").WithCode("test_code"),
			message: "test error",
			status:  codes.Unknown,
			details: map[string]string{"code": "test_code", "retryable": "false"},
		},
		{
			desc:    "error retryable",
			err:     errors.New("test error").Retryable(),
			message: "test error",
			status:  codes.Unknown,
			details: map[string]string{"code": "UNKNOWN", "retryable": "true"},
		},
		{
			desc:    "error with kind",
			err:     errors.New("test error").WithKind(errors.KindNotFound),
			message: "test error",
			status:  codes.NotFound,
			details: map[string]string{"code": "UNKNOWN", "retryable": "false"},
		},
		{
			desc:    "error with reasons",
			err:     errors.New("test error").WithCause(errors.New("reason 1")).WithCause(errors.New("reason 2")),
			message: "test error",
			status:  codes.Unknown,
			details: map[string]string{"code": "UNKNOWN", "retryable": "false", "reasons": "8#reason 18#reason 2"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			handler := func(ctx context.Context, req any) (any, error) {
				return nil, tc.err
			}

			inter := interceptor.UnaryErrorHandler()

			_, err := inter(context.Background(), nil, nil, handler)
			require.Error(t, err)

			st := status.Convert(err)

			assert.Equal(t, tc.status, st.Code())
			assert.Equal(t, tc.message, st.Message())

			details := st.Details()
			require.Len(t, details, 1)

			info, ok := details[0].(*errdetails.ErrorInfo)
			require.True(t, ok)

			assert.Equal(t, tc.details, info.Metadata)
		})
	}
}
