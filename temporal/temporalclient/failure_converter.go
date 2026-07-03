package temporalclient

import (
	"encoding/json"

	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/temporal"

	failurepb "go.temporal.io/api/failure/v1"

	"github.com/lcnascimento/go-kit/errors"
)

type failureConverter struct {
	converter *temporal.DefaultFailureConverter
}

func NewFailureConverter() converter.FailureConverter {
	return &failureConverter{
		converter: temporal.NewDefaultFailureConverter(temporal.DefaultFailureConverterOptions{
			EncodeCommonAttributes: false,
		}),
	}
}

func (f *failureConverter) ErrorToFailure(err error) *failurepb.Failure {
	kind := errors.Kind(err)
	code := errors.Code(err)
	reasons := errors.Reasons(err)
	retryable := errors.IsRetryable(err)

	isCustom := kind != errors.KindUnknown || code != errors.CodeUnknown || retryable
	if !isCustom {
		return f.converter.ErrorToFailure(err)
	}

	var cause error
	if len(reasons) > 0 {
		cause = errors.New("%s", reasons[len(reasons)-1])
	}

	err = temporal.NewApplicationErrorWithOptions(err.Error(), string(code), temporal.ApplicationErrorOptions{
		NonRetryable: !retryable,
		Cause:        cause,
		Details: []any{
			kind,
			reasons,
		},
	})

	return f.converter.ErrorToFailure(err)
}

func (f *failureConverter) FailureToError(failure *failurepb.Failure) error {
	info := failure.GetApplicationFailureInfo()
	if info == nil {
		return f.converter.FailureToError(failure)
	}

	err := errors.New("%s", failure.GetMessage())

	if info.Type != "" {
		err = err.WithCode(errors.CodeType(info.Type))
	}

	if !info.GetNonRetryable() {
		err = err.Retryable()
	}

	if info.Details == nil || len(info.Details.Payloads) == 0 {
		return err
	}

	var kind errors.KindType
	_ = json.Unmarshal(info.Details.Payloads[0].GetData(), &kind)

	if kind != "" && kind != errors.KindUnknown {
		err = err.WithKind(kind)
	}

	if len(info.Details.Payloads) == 1 {
		return err
	}

	var reasons []string
	_ = json.Unmarshal(info.Details.Payloads[1].GetData(), &reasons)

	for i := len(reasons) - 1; i >= 0; i-- {
		err = err.WithCause(errors.New("%s", reasons[i]))
	}

	return err
}
