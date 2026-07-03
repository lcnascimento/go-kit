package interceptor

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/log"
)

// UnaryErrorHandler returns a new unary interceptor suitable for request error handling.
func UnaryErrorHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = encodeError(onPanic(ctx, r))
			}
		}()

		resp, err = handler(ctx, req)
		if err != nil {
			err = encodeError(err)
		}

		return resp, err
	}
}

// StreamErrorHandler returns a new stream interceptor suitable for request error handling.
func StreamErrorHandler() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = encodeError(onPanic(ss.Context(), r))
			}
		}()

		err = handler(srv, ss)
		if err != nil {
			err = encodeError(err)
		}

		return err
	}
}

func encodeError(err error) error {
	kind := errors.Kind(err)
	reasons := errors.SafeReasons(err)

	details := &errdetails.ErrorInfo{
		Metadata: map[string]string{
			"code":      string(errors.Code(err)),
			"retryable": strconv.FormatBool(errors.IsRetryable(err)),
		},
	}

	if len(reasons) > 0 {
		details.Metadata["reasons"] = encodeReasons(reasons)
	}

	st, err := status.New(kindToGRPCStatusCode(kind), err.Error()).WithDetails(details)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return st.Err()
}

func encodeReasons(reasons []string) string {
	var encoded string

	for _, reason := range reasons {
		encoded += fmt.Sprintf("%d#%s", len(reason), reason)
	}

	return encoded
}

func kindToGRPCStatusCode(kind errors.KindType) codes.Code {
	switch kind {
	case errors.KindInvalidInput:
		return codes.InvalidArgument
	case errors.KindUnauthenticated:
		return codes.Unauthenticated
	case errors.KindUnauthorized:
		return codes.PermissionDenied
	case errors.KindNotFound:
		return codes.NotFound
	case errors.KindConflict:
		return codes.FailedPrecondition
	case errors.KindInternal, errors.KindCritical, errors.KindFatal:
		return codes.Internal
	case errors.KindResourceExhausted:
		return codes.ResourceExhausted
	default:
		return codes.Unknown
	}
}

func onPanic(ctx context.Context, cause any) error {
	logger.Critical(
		ctx,
		ErrPanic,
		log.Any("exception.message", cause),
		log.String("exception.stack", string(debug.Stack())),
	)

	return ErrPanic
}
