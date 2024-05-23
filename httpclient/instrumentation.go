package httpclient

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/log"
)

const (
	pkg         = "github.com/lcnascimento/go-kit/httpclient"
	counterName = "http_client_requests_count"
	latencyName = "http_client_requests_latency"
)

var (
	tracer          trace.Tracer
	requestsCounter metric.Int64Counter
	requestsLatency metric.Int64Histogram
)

func init() {
	ctx := context.Background()

	var err error

	tracer = otel.Tracer(pkg)
	meter := otel.Meter(pkg)

	requestsCounter, err = meter.Int64Counter(
		counterName,
		metric.WithDescription("Counts how many HTTP requests are made to an external system"),
		metric.WithUnit("1"),
	)
	if err != nil {
		log.Fatal(ctx, err)
	}

	requestsLatency, err = meter.Int64Histogram(
		latencyName,
		metric.WithDescription("Measures how many milliseconds HTTP requests are wasting"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatal(ctx, err)
	}
}

func (c *Client) onRequestStart(ctx context.Context, host, path, method string) trace.Span {
	ctx, span := tracer.Start(ctx, method, trace.WithAttributes())

	log.Debug(
		ctx,
		"http request started",
		log.String("host", host),
		log.String("path", path),
		log.String("method", method),
	)

	return span
}

func (c *Client) onRequestEnd(ctx context.Context, span trace.Span, host, path, method string, status int, start time.Time) {
	latency := time.Since(start)

	attrs := []attribute.KeyValue{
		attribute.String("host", host),
		attribute.String("path", path),
		attribute.String("method", method),
		attribute.Int("status_code", status),
	}

	span.SetAttributes(attrs...)

	mOption := metric.WithAttributeSet(attribute.NewSet(attrs...))
	requestsCounter.Add(ctx, 1, mOption)
	requestsLatency.Record(ctx, latency.Milliseconds(), mOption)

	log.Debug(
		ctx,
		"http request completed",
		log.String("host", host),
		log.String("path", path),
		log.String("method", method),
		log.Int("status_code", status),
		log.String("latency", latency.String()),
	)
}

func (c *Client) onParseURLError(ctx context.Context, url string, err error) error {
	err = ErrParseURL(err)
	log.Error(ctx, err, log.String("url", url))

	return err
}

func (c *Client) onBuildRequestError(ctx context.Context, err error) error {
	err = ErrBuildRequestError(err)
	log.Error(ctx, err)

	return err
}

func (c *Client) onRequestError(ctx context.Context, err error) error {
	err = ErrRequestError(err)
	log.Error(ctx, err)

	return err
}

func (c *Client) onBodyReadError(ctx context.Context, err error) error {
	err = ErrBodyReadError(err)
	log.Error(ctx, err)

	return err
}

func (c *Client) onUnexpectedStatusCode(ctx context.Context, code int, b []byte) error {
	var body map[string]any
	_ = json.Unmarshal(b, &body)

	attrs := []slog.Attr{
		slog.Int("status_code", code),
	}
	if len(body) > 0 {
		attrs = append(attrs, slog.Any("body", body))
	}

	err := ErrUnexpectedStatusCode(code)

	if code < http.StatusInternalServerError {
		log.Warn(ctx, err.Error(), attrs...)
	} else {
		log.Error(ctx, err, attrs...)
	}

	return err
}
