package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"

	URL "net/url"

	"github.com/lcnascimento/go-kit/errors"
)

const (
	defaultTimeoutInSeconds = 30
	contentType             = "application/json"
)

type httpClientProvider interface {
	Do(request *http.Request) (*http.Response, error)
}

type client struct {
	http    httpClientProvider
	timeout time.Duration
}

func New(opts ...Option) Client {
	client := &client{
		http:    &http.Client{},
		timeout: time.Second * defaultTimeoutInSeconds,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *client) Patch(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers["content-type"] = contentType

	return c.processRequest(ctx, "PATCH", request, opts...)
}

func (c *client) Put(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers["content-type"] = contentType

	return c.processRequest(ctx, "PUT", request, opts...)
}

func (c *client) Post(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers["content-type"] = contentType

	return c.processRequest(ctx, "POST", request, opts...)
}

func (c *client) PostForm(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers["content-type"] = "application/x-www-form-urlencoded"

	return c.processRequest(ctx, "POST", request, opts...)
}

func (c *client) Delete(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	return c.processRequest(ctx, "DELETE", request, opts...)
}

func (c *client) Get(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers["content-type"] = contentType

	return c.processRequest(ctx, "GET", request, opts...)
}

func (c *client) processRequest(ctx context.Context, method string, request *Request, opts ...RequestOption) (rst Result, err error) {
	defer func() {
		if err != nil && ctx.Err() != nil {
			rst = Result{}
			err = errors.ErrContextCanceled
		}
	}()

	for _, opt := range opts {
		opt(request)
	}

	queryValues := URL.Values{}

	for key, value := range request.QueryParams {
		queryValues.Add(key, value)
	}

	uri := request.Host + request.Path
	for p, v := range request.PathParams {
		if strings.Contains(uri, ":"+p) {
			uri = strings.ReplaceAll(uri, ":"+p, v)
		}
	}

	url, err := URL.Parse(uri)
	if err != nil {
		return Result{}, errors.New("error on parsing the request url")
	}

	url.RawQuery = queryValues.Encode()

	httpRequest, err := http.NewRequestWithContext(ctx, method, url.String(), bytes.NewBuffer(request.Body))
	if err != nil {
		return Result{}, err
	}

	for key, value := range request.Headers {
		httpRequest.Header.Add(key, value)
	}

	ctx, span, start := c.onRequestStart(ctx, request.Host, request.Path, method)
	defer span.End()

	return c.doRequest(ctx, httpRequest, request, method, start, span)
}

func (c *client) doRequest(
	ctx context.Context, httpRequest *http.Request, request *Request,
	method string, start time.Time, span trace.Span,
) (Result, error) {
	res, err := c.http.Do(httpRequest)
	if err != nil {
		return Result{}, errors.ErrRequestError.WithCause(err).WithAttribute("http.request.body", string(request.Body))
	}

	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Result{}, err
	}

	c.onRequestEnd(ctx, request.Host, request.Path, method, res.StatusCode, start, span)

	result := Result{
		Response:   body,
		StatusCode: res.StatusCode,
	}

	configAcceptable := request.acceptStatusCodes != nil && request.acceptStatusCodes[res.StatusCode]
	defaultAcceptable := len(request.acceptStatusCodes) == 0 && res.StatusCode < http.StatusBadRequest

	if configAcceptable || defaultAcceptable {
		return result, nil
	}

	dErr := ErrUnexpectedStatusCode.WithKind(statusCodeToKind(res.StatusCode))
	if res.StatusCode >= http.StatusInternalServerError {
		dErr = dErr.Retryable()
	}

	dErr = dErr.
		WithAttribute("http.response.status_code", strconv.Itoa(res.StatusCode)).
		WithAttribute("http.response.body", string(body))

	return result, dErr
}
