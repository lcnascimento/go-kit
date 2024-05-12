package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"

	"github.com/lcnascimento/go-kit/trace"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("example")
}

func main() {
	_, span := tracer.Start(context.Background(), "main")
	defer span.End()

	fmt.Println("TraceID:", span.SpanContext().TraceID().String())
	fmt.Println("SpanID:", span.SpanContext().SpanID().String())
}
