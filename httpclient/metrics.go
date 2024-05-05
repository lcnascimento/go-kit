package httpclient

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	metricRequestsCounter         = "http_client_requests_count"
	metricRequestLatencyHistogram = "http_client_requests_latency"
)

func (c *Client) mustInitMetrics() {
	var err error

	c.metricRequestsCount, err = c.meter.Int64Counter(
		metricRequestsCounter,
		metric.WithDescription("Counts how many HTTP requests are made to an external system"),
		metric.WithUnit("1"),
	)
	if err != nil {
		panic(ErrMetricInitialization(metricRequestsCounter))
	}

	c.metricRequestsLatency, err = c.meter.Int64Histogram(
		metricRequestLatencyHistogram,
		metric.WithDescription("Measures how many milliseconds HTTP requests are wasting"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		panic(ErrMetricInitialization(metricRequestLatencyHistogram))
	}
}

func (c *Client) onRequestEnd(ctx context.Context, host, path, method string, status int, start time.Time) {
	latency := time.Since(start).Milliseconds()

	attrs := []attribute.KeyValue{
		attribute.String("host", host),
		attribute.String("path", path),
		attribute.String("method", method),
		attribute.Int("status_code", status),
	}

	c.metricRequestsCount.Add(ctx, 1, metric.WithAttributeSet(attribute.NewSet(attrs...)))
	c.metricRequestsLatency.Record(ctx, latency, metric.WithAttributeSet(attribute.NewSet(attrs...)))
}
