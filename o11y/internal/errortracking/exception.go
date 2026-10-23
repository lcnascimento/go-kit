package errortracking

import (
	"context"
	"time"

	"github.com/posthog/posthog-go"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

const (
	memberUserID         = "user.id"
	memberOrganizationID = "organization.id"
	memberSessionID      = "posthog.session_id"
)

func newException(ctx context.Context, err error) posthog.Exception {
	bag := baggage.FromContext(ctx)
	extractor := stackTraceExtractor()

	return posthog.Exception{
		DistinctId: distinctID(bag),
		Timestamp:  time.Now(),
		Properties: properties(ctx, bag, err),
		ExceptionList: []posthog.ExceptionItem{{
			Type:       string(errors.Code(err)),
			Value:      err.Error(),
			Stacktrace: stackTrace(extractor),
		}},
		DebugImages: extractor.GetDebugImages(),
	}
}

func distinctID(bag baggage.Baggage) string {
	if userID := bag.Member(memberUserID).Value(); userID != "" {
		return userID
	}

	return global.Config().ServiceName
}

// properties merges, from lowest to highest precedence, the baggage, the error's attributes and the fixed properties.
func properties(ctx context.Context, bag baggage.Baggage, err error) posthog.Properties {
	props := posthog.NewProperties()

	for _, member := range bag.Members() {
		if member.Key() != memberSessionID {
			props.Set(member.Key(), member.Value())
		}
	}

	for key, value := range errors.Attributes(err) {
		props.Set(key, value)
	}
	if reasons := errors.Reasons(err); len(reasons) > 0 {
		props.Set("reasons", reasons)
	}

	cfg := global.Config()
	props.
		Set("kind", string(errors.Kind(err))).
		Set("code", string(errors.Code(err))).
		Set("severity", errors.Severity(err).String()).
		Set("retryable", errors.IsRetryable(err)).
		Set("service.name", cfg.ServiceName)

	if environment := cfg.ResourceAttributes[semconv.DeploymentEnvironmentKey]; environment != "" {
		props.Set("deployment.environment", environment)
	}

	if spanContext := trace.SpanContextFromContext(ctx); spanContext.IsValid() {
		props.Set("$trace_id", spanContext.TraceID().String()).Set("$span_id", spanContext.SpanID().String())
	}

	if sessionID := bag.Member(memberSessionID).Value(); sessionID != "" {
		props.Set("$session_id", sessionID)
	}

	if organizationID := bag.Member(memberOrganizationID).Value(); organizationID != "" {
		props.Set("$groups", posthog.Groups{"organization": organizationID})
	}

	if bag.Member(memberUserID).Value() == "" {
		props.Set("$process_person_profile", false)
	}

	return props
}
