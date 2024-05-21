package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	URL "net/url"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const defaultTimeoutInSeconds = 30

type httpClientProvider interface {
	Do(request *http.Request) (*http.Response, error)
}

// Client provides methods for making REST requests.
type Client struct {
	http    httpClientProvider
	timeout time.Duration
}

// New creates a new Client instance.
func New(opts ...Option) *Client {
	client := &Client{
		http: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
		timeout: time.Second * defaultTimeoutInSeconds,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Patch execute a http PATCH method with application/json headers.
func (c *Client) Patch(ctx context.Context, req *Request, opts ...RequestOption) (*Response, error) {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	req.Headers["content-type"] = "application/json"

	return c.processRequest(ctx, "PATCH", req, opts...)
}

// Put execute a http PUT method with application/json headers.
func (c *Client) Put(ctx context.Context, req *Request, opts ...RequestOption) (*Response, error) {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	req.Headers["content-type"] = "application/json"

	return c.processRequest(ctx, "PUT", req, opts...)
}

// Post execute a http POST method with application/json headers.
func (c *Client) Post(ctx context.Context, req *Request, opts ...RequestOption) (*Response, error) {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	req.Headers["content-type"] = "application/json"

	return c.processRequest(ctx, "POST", req, opts...)
}

// PostForm execute a http POST method with x-www-form-urlencoded headers.
func (c *Client) PostForm(ctx context.Context, req *Request, opts ...RequestOption) (*Response, error) {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	req.Headers["content-type"] = "application/x-www-form-urlencoded"

	return c.processRequest(ctx, "POST", req, opts...)
}

// Delete execute a http DELETE method with application/json headers.
func (c *Client) Delete(ctx context.Context, req *Request, opts ...RequestOption) (*Response, error) {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	return c.processRequest(ctx, "DELETE", req, opts...)
}

// Get execute a http GET method.
func (c *Client) Get(ctx context.Context, req *Request, opts ...RequestOption) (*Response, error) {
	return c.processRequest(ctx, "GET", req, opts...)
}

func (c *Client) processRequest(ctx context.Context, method string, req *Request, opts ...RequestOption) (*Response, error) {
	for _, opt := range opts {
		opt(req)
	}

	url, err := c.buildURL(ctx, req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, url.String(), bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, c.onBuildRequestError(ctx, err)
	}

	for key, value := range req.Headers {
		request.Header.Add(key, value)
	}

	start := time.Now()

	span := c.onRequestStart(ctx, request.Host, req.Path, method)

	res, err := c.http.Do(request)
	if err != nil {
		return nil, c.onRequestError(ctx, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, c.onBodyReadError(ctx, err)
	}

	c.onRequestEnd(ctx, span, request.Host, req.Path, method, res.StatusCode, start)

	response := &Response{
		Body:       body,
		StatusCode: res.StatusCode,
	}

	configNonAcceptable := req.acceptableStatusCodes != nil && !req.acceptableStatusCodes[res.StatusCode]
	defaultNonAcceptable := res.StatusCode >= http.StatusBadRequest

	if configNonAcceptable || defaultNonAcceptable {
		return response, c.onUnexpectedStatusCode(ctx, response.StatusCode, response.Body)
	}

	return response, nil
}

func (c *Client) buildURL(ctx context.Context, req *Request) (*URL.URL, error) {
	queryValues := URL.Values{}
	for key, value := range req.QueryParams {
		queryValues.Add(key, value)
	}

	uri := req.Host + req.Path
	for p, v := range req.PathParams {
		if strings.Contains(uri, ":"+p) {
			uri = strings.ReplaceAll(uri, ":"+p, v)
		}
	}

	url, err := URL.Parse(uri)
	if err != nil {
		return nil, c.onParseURLError(ctx, uri, err)
	}

	url.RawQuery = queryValues.Encode()

	return url, nil
}
