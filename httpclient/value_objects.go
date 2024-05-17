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

// HTTPHeaders is a map containing the relation key=value of the headers used on the http rest request.
type HTTPHeaders map[string]string

// HTTPPathParams is a map containing the relation key=value of the path params used on the http rest request.
// It will be used to replace values given in Path parameter.
type HTTPPathParams map[string]string

// HTTPQueryParams is a map containing the relation key=value of the query params used on the http rest request.
type HTTPQueryParams map[string]string

// HTTPRequest are the params used to build a new http rest request.
type HTTPRequest struct {
	Host        string
	Path        string
	Body        []byte
	Headers     HTTPHeaders
	QueryParams HTTPQueryParams
	PathParams  HTTPPathParams
}

// HTTPResult are the params returned from the client HTTP request.
type HTTPResult struct {
	StatusCode int
	Response   []byte
}
