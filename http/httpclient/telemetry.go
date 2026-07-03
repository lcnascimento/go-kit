package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.37.0/httpconv"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	"github.com/lcnascimento/go-kit/o11y/metric"
)

var (
	pkg = "github.com/lcnascimento/go-kit/http/httpclient"

	meter  = otel.Meter(pkg)
	tracer = otel.Tracer(pkg)

	activeRequestsMetric, _  = httpconv.NewClientActiveRequests(meter)
	requestDurationMetric, _ = httpconv.NewClientRequestDuration(meter)
	requestsCounter          = metric.MustIntCounter(meter, "http.client.request.total", "Number of HTTP requests")
)

func (c *client) onRequestStart(ctx context.Context, host, path, method string) (context.Context, trace.Span, time.Time) {
	attrs := []attribute.KeyValue{
		attribute.String(string(semconv.NetworkPeerAddressKey), host),
		attribute.String(string(semconv.HTTPRequestMethodKey), method),
		attribute.String(string(semconv.HTTPRouteKey), path),
	}

	operation := fmt.Sprintf("%s %s", method, path)
	ctx, span := tracer.Start(ctx, operation, trace.WithAttributes(attrs...))

	const port = 80

	activeRequestsMetric.Add(ctx, 1, host, port, attrs...)

	return ctx, span, time.Now()
}

func (c *client) onRequestEnd(ctx context.Context, host, path, method string, status int, start time.Time, span trace.Span) {
	latency := time.Since(start).Milliseconds()

	attrs := []attribute.KeyValue{
		attribute.String(string(semconv.NetworkPeerAddressKey), host),
		attribute.String(string(semconv.HTTPRequestMethodKey), method),
		attribute.String(string(semconv.HTTPRouteKey), path),
	}

	const port = 80

	activeRequestsMetric.Add(ctx, -1, host, port, attrs...)

	attrs = append(attrs, attribute.Int(string(semconv.HTTPResponseStatusCodeKey), status))
	requestsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	requestDurationMetric.Record(ctx, float64(latency), httpconv.RequestMethodAttr(method), host, port, attrs...)

	if status >= http.StatusInternalServerError {
		span.SetStatus(codes.Error, semconv.HTTPResponseStatusCode(status).Value.AsString())
	} else {
		span.SetStatus(codes.Unset, "")
	}
}
