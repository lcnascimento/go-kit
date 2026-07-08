package errors

var (
	// ErrNotImplemented indicates that a given feature is not implemented yet.
	ErrNotImplemented error = New("feature not implemented yet").
				WithCode("FEATURE_NOT_IMPLEMENTED").
				WithKind(KindInternal)

	// ErrMock is a fake mocked that should be used in test scenarios.
	ErrMock error = New("mocked error").
		WithCode("MOCKED_ERROR").
		WithKind(KindInternal)

	// ErrResourceNotFound indicates that a desired resource was not found.
	ErrResourceNotFound error = New("resource not found").
				WithCode("RESOURCE_NOT_FOUND").
				WithKind(KindNotFound)

	// ErrRequestUnauthenticated indicates that the request made to an API has missing or invalid credentials.
	ErrRequestUnauthenticated = New("unauthenticated request").
					WithCode("ERR_UNAUTHENTICATED").
					WithKind(KindUnauthenticated)

	// ErrRequestUnauthorized indicates that the request made to an API has missing or invalid credentials.
	ErrRequestUnauthorized = New("unauthorized request").
				WithCode("ERR_UNAUTHORIZED").
				WithKind(KindUnauthorized)

	// ErrRequestError occurs when there is any error on requests to external systems.
	ErrRequestError = New("request to external system failed").
			WithCode("ERR_REQUEST_TO_EXTERNAL_SYSTEM_FAILED").
			WithKind(KindInternal).
			Retryable()

	// ErrUnexpectedResponseStatus occurs when a request to an external system returns an unexpected status code.
	ErrUnexpectedResponseStatus = New("unexpected status response from external system").
					WithCode("ERR_UNEXPECTED_RESPONSE_STATUS").
					WithKind(KindInternal).
					Retryable()

	// ErrCastPayload indicates an issue during response payload casting.
	ErrCastPayload = New("could not cast payload").
			WithCode("ERR_CAST_PAYLOAD").
			WithKind(KindInternal)

	// ErrInvalidInput occurs when a given input is invalid.
	ErrInvalidInput = New("invalid input").
			WithCode("ERR_INVALID_INPUT").
			WithKind(KindInvalidInput)

	// ErrInvalidAccountID occurs when a given account ID is invalid.
	ErrInvalidAccountID = New("invalid account ID").
				WithCode("ERR_INVALID_ACCOUNT_ID").
				WithKind(KindInvalidInput)

	// ErrFetchFeatureFlag occurs when a request made to a Feature Flag service fails.
	ErrFetchFeatureFlag = New("could not fetch feature flag").
				WithCode("ERR_FETCH_FEATURE_FLAG").
				WithKind(KindInternal).
				Retryable()

	// ErrContextCanceled indicates that an operation was canceled, typically by a context cancellation.
	ErrContextCanceled = New("operation canceled").
				WithCode("ERR_CONTEXT_CANCELED").
				WithKind(KindCanceled).
				Retryable()
)

type (
	CodeType string
	KindType string
)

const CodeUnknown CodeType = "UNKNOWN"

const (
	KindUnknown            KindType = "UNKNOWN"
	KindConflict           KindType = "CONFLICT"
	KindInternal           KindType = "INTERNAL"
	KindInvalidInput       KindType = "INVALID_INPUT"
	KindNotFound           KindType = "NOT_FOUND"
	KindUnauthenticated    KindType = "UNAUTHENTICATED"
	KindUnauthorized       KindType = "UNAUTHORIZED"
	KindUnprocessable      KindType = "UNPROCESSABLE"
	KindResourceExhausted  KindType = "RESOURCE_EXHAUSTED"
	KindServiceUnavailable KindType = "SERVICE_UNAVAILABLE"
	KindCritical           KindType = "CRITICAL"
	KindFatal              KindType = "FATAL"
	KindCanceled           KindType = "CANCELED"
	KindWarn               KindType = "WARN"
)

type SeverityType int

const (
	SeverityWarn SeverityType = iota

	// SeverityError indicates a regular error.
	SeverityError

	// SeverityCritical indicates a critical error.
	// Typically, this kind of error is not expected to happen and may require immediate attention.
	SeverityCritical

	// SeverityFatal indicates a fatal error.
	// Typically, this kind of error is not expected to happen and will cause the application to crash.
	SeverityFatal
)

// String returns the string representation of the SeverityType.
func (s SeverityType) String() string {
	return severityNames[s]
}

var severityNames = []string{
	"WARN",
	"ERROR",
	"CRITICAL",
	"FATAL",
}

// AttributeSet is a map of string key-value pairs that can be used to add additional information to an error.
type AttributeSet map[string]string

// Merge merges the given AttributeSet into the current one.
func (s AttributeSet) Merge(other AttributeSet) {
	for key, value := range other {
		s[key] = value
	}
}
