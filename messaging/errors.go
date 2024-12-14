package messaging

import "github.com/lcnascimento/go-kit/errors"

var (
	// ErrCreateBaggage is returned when failed to create baggage.
	ErrCreateBaggage errors.CustomError = errors.New("failed to create baggage").
				WithCode("CREATE_BAGGAGE_ERROR").
				WithKind(errors.KindInternal)

	// ErrCreateBaggageMember is returned when failed to create baggage member.
	ErrCreateBaggageMember errors.CustomError = errors.New("failed to create baggage member").
				WithCode("CREATE_BAGGAGE_MEMBER_ERROR").
				WithKind(errors.KindInternal)

	// ErrSetBaggageMember is returned when failed to set baggage member.
	ErrSetBaggageMember errors.CustomError = errors.New("failed to set baggage member").
				WithCode("SET_BAGGAGE_MEMBER_ERROR").
				WithKind(errors.KindInternal)
)
