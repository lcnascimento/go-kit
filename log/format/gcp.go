package format

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/errors"
)

// GCPCloudLoggingLogFormatterParams encapsulates necessary parameters to construct a GCP Cloud Logging LogFormatter.
type GCPCloudLoggingLogFormatterParams struct {
	ProjectID          string
	ApplicationName    string
	ApplicationVersion string
}

type gcpCloudLoggingLogFormatter struct {
	projectID          string
	applicationName    string
	applicationVersion string
}

// NewGCPCloudLogging creates a new GCP Cloud Logging LogFormatter.
func NewGCPCloudLogging(params GCPCloudLoggingLogFormatterParams) (*gcpCloudLoggingLogFormatter, error) {
	if params.ApplicationName == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if params.ApplicationVersion == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationVersion")
	}

	if params.ProjectID == "" {
		return nil, errors.NewMissingRequiredDependency("ProjectID")
	}

	return &gcpCloudLoggingLogFormatter{
		projectID:          params.ProjectID,
		applicationName:    params.ApplicationName,
		applicationVersion: params.ApplicationVersion,
	}, nil
}

// MustNewGCPCloudLogging creates a new GCP Cloud Logging LogFormatter.
// It panics if any error is found.
func MustNewGCPCloudLogging(params GCPCloudLoggingLogFormatterParams) *gcpCloudLoggingLogFormatter {
	formatter, err := NewGCPCloudLogging(params)
	if err != nil {
		panic(err)
	}

	return formatter
}

// Format formats the log payload that will be rendered in accordance with Cloud Logging standards..
func (b gcpCloudLoggingLogFormatter) Format(ctx context.Context, in LogInput) any {
	payload := map[string]any{
		"severity": in.Level,
		"time":     in.Timestamp.Format(time.RFC3339),
		"message":  in.Message,
	}

	if in.Payload != nil {
		payload["payload"] = in.Payload
	}

	contextKeys := extractContextKeysFromContext(ctx, in.ContextKeys)
	if len(contextKeys) > 0 {
		payload["context"] = contextKeys
	}

	if len(in.Attributes) > 0 {
		payload["attributes"] = in.Attributes
	}

	if isError(in.Level) {
		// Necessary to link error to Cloud Error Reporting.
		// More details in: https://cloud.google.com/error-reporting/docs/formatting-error-messages
		payload["@type"] = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
		payload["serviceContext"] = map[string]interface{}{
			"service": b.applicationName,
			"version": b.applicationVersion,
		}
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().TraceID().IsValid() {
		return payload
	}

	if isError(in.Level) {
		span.SetStatus(codes.Error, in.Message)
	}

	// Necessary to link with Cloud Trace
	// More details in: https://cloud.google.com/logging/docs/structured-logging
	payload["logging.googleapis.com/trace"] = fmt.Sprintf("projects/%s/traces/%s", b.projectID, span.SpanContext().TraceID().String())
	payload["logging.googleapis.com/spanId"] = span.SpanContext().SpanID().String()
	payload["logging.googleapis.com/trace_sampled"] = span.SpanContext().IsSampled()

	return payload
}
