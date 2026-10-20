package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.37.0/httpconv"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	"github.com/lcnascimento/go-kit/http/httpserver/middlewares/internal"
	"github.com/lcnascimento/go-kit/o11y/baggage"
	"github.com/lcnascimento/go-kit/o11y/log"
	"github.com/lcnascimento/go-kit/o11y/metric"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/http/httpserver/middlewares"
	logger = log.MustNewLogger(pkg)
	meter  = otel.Meter(pkg)
	tracer = otel.Tracer(pkg)
)

var (
	totalRequestsMetric      = metric.MustIntCounter(meter, "http.server.request.total", "Total number of HTTP Requests made to the server")
	requestSizeMetric, _     = httpconv.NewServerRequestBodySize(meter)
	requestDurationMetric, _ = httpconv.NewServerRequestDuration(
		meter,
		//nolint:mnd // Histogram bucket boundaries, in seconds.
		otelmetric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10),
	)
)

func Telemetry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		start := time.Now()

		rww := internal.NewRespWriterWrapper(w, func(int64) {})

		// Wrap w to use our ResponseWriter methods while also exposing
		// other interfaces that w may implement (http.CloseNotifier,
		// http.Flusher, http.Hijacker, http.Pusher, io.ReaderFrom).
		w = httpsnoop.Wrap(w, httpsnoop.Hooks{
			Header: func(httpsnoop.HeaderFunc) httpsnoop.HeaderFunc {
				return rww.Header
			},
			Write: func(httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				return rww.Write
			},
			WriteHeader: func(httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				return rww.WriteHeader
			},
			Flush: func(httpsnoop.FlushFunc) httpsnoop.FlushFunc {
				return rww.Flush
			},
		})

		var path string
		if route := mux.CurrentRoute(r); route != nil {
			path, _ = route.GetPathTemplate()
		}

		operation := spanName(r.Method, path)

		ctx, span := tracer.Start(ctx, operation, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		status := rww.StatusCode()
		duration := time.Since(start)

		logRequest(ctx, r, operation, path, status)
		measureRequest(ctx, r, path, status, duration)
		trackRequest(ctx, r, path, status, span)
	})
}

func spanName(method, pathTpl string) string {
	if pathTpl == "" {
		return method
	}

	return fmt.Sprintf("%s %s", method, pathTpl)
}

func logRequest(ctx context.Context, r *http.Request, operation, pathTpl string, status int) {
	logger.Debug(
		ctx, operation,
		log.String(string(semconv.NetworkPeerAddressKey), r.Host),
		log.String(string(semconv.HTTPRequestMethodKey), r.Method),
		log.String(string(semconv.HTTPRouteKey), pathTpl),
		log.String(string(semconv.URLPathKey), r.URL.Path),
		log.String(string(semconv.UserAgentOriginalKey), r.UserAgent()),
		log.Int(string(semconv.HTTPResponseStatusCodeKey), status),
	)
}

func measureRequest(ctx context.Context, r *http.Request, pathTpl string, status int, duration time.Duration) {
	method := normalizeMethod(r.Method)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	attrs := []attribute.KeyValue{
		attribute.String(string(semconv.HTTPRouteKey), pathTpl),
		attribute.Int(string(semconv.HTTPResponseStatusCodeKey), status),
	}

	counterAttrs := append([]attribute.KeyValue{
		attribute.String(string(semconv.HTTPRequestMethodKey), string(method)),
	}, attrs...)

	totalRequestsMetric.Add(ctx, 1, metric.WithAttributes(counterAttrs...))
	requestDurationMetric.Record(ctx, duration.Seconds(), method, scheme, attrs...)

	if r.ContentLength >= 0 {
		requestSizeMetric.Record(ctx, r.ContentLength, method, scheme, attrs...)
	}
}

func normalizeMethod(method string) httpconv.RequestMethodAttr {
	switch m := httpconv.RequestMethodAttr(method); m {
	case httpconv.RequestMethodConnect, httpconv.RequestMethodDelete, httpconv.RequestMethodGet,
		httpconv.RequestMethodHead, httpconv.RequestMethodOptions, httpconv.RequestMethodPatch,
		httpconv.RequestMethodPost, httpconv.RequestMethodPut, httpconv.RequestMethodTrace:
		return m
	default:
		return httpconv.RequestMethodOther
	}
}

func trackRequest(ctx context.Context, r *http.Request, pathTpl string, status int, span trace.Span) {
	bag := baggage.FromContext(ctx)

	attrs := []attribute.KeyValue{ //nolint:prealloc // Baggage is small; a hand-counted capacity is not worth it.
		attribute.String(string(semconv.HTTPRequestMethodKey), r.Method),
		attribute.String(string(semconv.HTTPRouteKey), pathTpl),
		attribute.String(string(semconv.URLPathKey), r.URL.Path),
		attribute.Int(string(semconv.HTTPResponseStatusCodeKey), status),
	}

	for _, member := range bag.Members() {
		attrs = append(attrs, attribute.String("bag."+member.Key(), member.Value()))
	}

	if status >= http.StatusInternalServerError {
		span.SetStatus(codes.Error, semconv.HTTPResponseStatusCode(status).Value.AsString())
	} else {
		span.SetStatus(codes.Unset, "")
	}

	span.SetAttributes(attrs...)
}
