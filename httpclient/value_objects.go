package httpclient

import (
	"time"
)

// Option is a type to set HTTP Client options.
type Option func(*Client)

// WithTimeout instructs the HTTP Client to cancel any requests that exceeds the given timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// Headers is a map containing the relation key=value of the headers used on the http rest request.
type Headers map[string]string

// PathParams is a map containing the relation key=value of the path params used on the http rest request.
// It will be used to replace values given in Path parameter.
type PathParams map[string]string

// QueryParams is a map containing the relation key=value of the query params used on the http rest request.
type QueryParams map[string]string

// Request are the params used to build a new http rest request.
type Request struct {
	Host        string
	Path        string
	Body        []byte
	Headers     Headers
	QueryParams QueryParams
	PathParams  PathParams

	acceptableStatusCodes map[int]bool
}

// RequestOption is a type to set HTTP Request options.
type RequestOption func(*Request)

// WithAcceptStatusCode indicates that the current request accepts the given status code.
// It is possible to use this option multiple times for different status codes in the same request.
func WithAcceptStatusCode(code int) RequestOption {
	return func(r *Request) {
		if r.acceptableStatusCodes == nil {
			r.acceptableStatusCodes = map[int]bool{}
		}

		r.acceptableStatusCodes[code] = true
	}
}

// Response encapsulates data returned from the client HTTP request.
type Response struct {
	StatusCode int
	Body       []byte
}
