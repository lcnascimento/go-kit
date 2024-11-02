package httpclient

import "context"

// HTTPClient is responsible for performing http requests.
type HTTPClient interface {
	// Get execute a http GET method.
	Get(context.Context, *Request, ...RequestOption) (*Response, error)

	// Post execute a http POST method with application/json headers.
	Post(context.Context, *Request, ...RequestOption) (*Response, error)

	// PostForm execute a http POST method with application/x-www-form-urlencoded headers.
	PostForm(context.Context, *Request, ...RequestOption) (*Response, error)

	// Put execute a http PUT method with application/json headers.
	Put(context.Context, *Request, ...RequestOption) (*Response, error)

	// Patch execute a http PATCH method with application/json headers.
	Patch(context.Context, *Request, ...RequestOption) (*Response, error)

	// Delete execute a http DELETE method with application/json headers.
	Delete(context.Context, *Request, ...RequestOption) (*Response, error)
}
