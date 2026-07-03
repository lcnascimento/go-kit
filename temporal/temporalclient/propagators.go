package temporalclient

import (
	"go.temporal.io/sdk/workflow"

	"github.com/lcnascimento/go-kit/temporal/internal/propagators"
)

func Propagators() []workflow.ContextPropagator {
	return []workflow.ContextPropagator{
		propagators.NewBaggagePropagator(),
	}
}
