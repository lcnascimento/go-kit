package errors

var (
	// ErrResourceNotFound indicates that a desired resource was not found.
	ErrResourceNotFound error = New("resource not found").WithKind(KindNotFound).WithCode("RESOURCE_NOT_FOUND")

	// ErrNotImplemented indicates that a given feature is not implemented yet.
	ErrNotImplemented error = New("feature not implemented yet").WithCode("FEATURE_NOT_IMPLEMENTED")

	// ErrMock is a fake mocked that should be used in test scenarios.
	ErrMock error = New("mocked error").WithCode("MOCKED_ERROR")
)

// CodeType is a string that contains error's code description.
type CodeType string

const (
	// CodeUnknown is the default code returned when the application doesn't attach any code into the error.
	CodeUnknown CodeType = "UNKNOWN"
)

// KindType is a string that contains error's kind description.
type KindType string

const (
	// KindUnknown is the default kind returned when the application doesn't attach any kind into the error.
	KindUnknown KindType = "UNKNOWN"

	// KindConflict are errors caused by requests with data that conflicts with the current state of the system.
	KindConflict KindType = "CONFLICT"

	// KindInternal are errors caused by some internal fail like failed IO calls or invalid memory states.
	KindInternal KindType = "INTERNAL"

	// KindInvalidInput are errors caused by some invalid values on the input.
	KindInvalidInput KindType = "INVALID_INPUT"

	// KindNotFound are errors caused by any required resources that not exists on the data repository.
	KindNotFound KindType = "NOT_FOUND"

	// KindUnauthenticated are errors caused by an unauthenticated call.
	KindUnauthenticated KindType = "UNAUTHENTICATED"

	// KindUnauthorized are errors caused by an unauthorized call.
	KindUnauthorized KindType = "UNAUTHORIZED"

	// KindResourceExhausted indicates some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
	KindResourceExhausted KindType = "RESOURCE_EXHAUSTED"
)
