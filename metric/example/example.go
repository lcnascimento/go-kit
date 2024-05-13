package main

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/lcnascimento/go-kit/metric"
)

var meter metric.Meter

func init() {
	meter = otel.Meter("example")
}

func main() {
	ctx := context.Background()

	counter, err := meter.Int64Counter("example_counter")
	if err != nil {
		panic(err)
	}

	counter.Add(ctx, 1)
}
