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

	"github.com/lcnascimento/go-kit/errors"
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
func (c *Client) Patch(ctx context.Context, request *HTTPRequest) (rst HTTPResult, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers["content-type"] = "application/json"

	return c.processRequest(ctx, "PATCH", request)
}

// Put execute a http PUT method with application/json headers.
func (c *Client) Put(ctx context.Context, request *HTTPRequest) (rst HTTPResult, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers["content-type"] = "application/json"

	return c.processRequest(ctx, "PUT", request)
}

// Post execute a http POST method with application/json headers.
func (c *Client) Post(ctx context.Context, request *HTTPRequest) (HTTPResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers["content-type"] = "application/json"

	return c.processRequest(ctx, "POST", request)
}

// PostForm execute a http POST method with x-www-form-urlencoded headers.
func (c *Client) PostForm(ctx context.Context, request *HTTPRequest) (HTTPResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers["content-type"] = "application/x-www-form-urlencoded"

	return c.processRequest(ctx, "POST", request)
}

// Delete execute a http DELETE method with application/json headers.
func (c *Client) Delete(ctx context.Context, request *HTTPRequest) (HTTPResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	return c.processRequest(ctx, "DELETE", request)
}

// Get execute a http GET method.
func (c *Client) Get(ctx context.Context, request *HTTPRequest) (HTTPResult, error) {
	return c.processRequest(ctx, "GET", request)
}

func (c *Client) processRequest(ctx context.Context, method string, request *HTTPRequest) (HTTPResult, error) {
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
		return HTTPResult{}, errors.New("error on parsing the request url")
	}

	url.RawQuery = queryValues.Encode()

	httpRequest, err := http.NewRequestWithContext(ctx, method, url.String(), bytes.NewBuffer(request.Body))
	if err != nil {
		return HTTPResult{}, err
	}

	for key, value := range request.Headers {
		httpRequest.Header.Add(key, value)
	}

	start := time.Now()

	span := c.onRequestStart(ctx, request.Host, request.Path, method)

	res, err := c.http.Do(httpRequest)
	if err != nil {
		return HTTPResult{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return HTTPResult{}, err
	}

	c.onRequestEnd(ctx, span, request.Host, request.Path, method, res.StatusCode, start)

	return HTTPResult{
		Response:   body,
		StatusCode: res.StatusCode,
	}, nil
}
