package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lcnascimento/go-kit/errors"
)

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
	case errors.KindInternal:
		return codes.Internal
	case errors.KindResourceExhausted:
		return codes.ResourceExhausted
	default:
		return codes.Unknown
	}
}

// ErrorHandlingUnaryServerInterceptor returns a new unary interceptor suitable for request error handling.
func ErrorHandlingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		resp, err = handler(ctx, req)
		if err != nil {
			err = status.Error(kindToGRPCStatusCode(errors.Kind(err)), err.Error())
		}

		return
	}
}

// ErrorHandlingStreamServerInterceptor returns a new stream interceptor suitable for request error handling.
func ErrorHandlingStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {

		err = handler(srv, ss)
		if err != nil {
			err = status.Error(kindToGRPCStatusCode(errors.Kind(err)), err.Error())
		}

		return
	}
}
