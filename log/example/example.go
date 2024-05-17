package main

import (
	"context"
	e "errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/log"
	"github.com/lcnascimento/go-kit/propagation"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("example")
}

func main() {
	foo := propagation.ContextKey("foo")
	bar := propagation.ContextKey("bar")

	log.SetLevel(log.LevelDebug)
	log.SetContextKeySet(propagation.ContextKeySet{
		foo: true,
		bar: true,
	})

	ctx := context.Background()
	ctx = context.WithValue(ctx, foo, "foo")
	ctx = context.WithValue(ctx, bar, "bar")

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	attr1 := log.String("attr1", "value1")
	attr2 := log.String("attr2", "value2")

	log.Debug(ctx, "Debug", attr1, attr2)
	log.Info(ctx, "Info", attr1, attr2)
	log.Warn(ctx, "Warn", attr1, attr2)
	log.Errorw(ctx, "Error", attr1, attr2)
	log.Criticalw(ctx, "Critical", attr1, attr2)

	log.Error(ctx, errDefault, attr1, attr2)
	log.Critical(ctx, errCritical, attr1, attr2)
}

var (
	errDefault = errors.New("default error").
			WithKind(errors.KindInvalidInput).
			WithCode("ERR_INVALID_INPUT").
			WithRootError(e.New("root default error"))

	errCritical = errors.New("critical error").
			WithKind(errors.KindUnexpected).
			WithCode("ERR_CRITICAL").
			WithRootError(e.New("root critical error")).
			WithStack().
			Retryable(true)
)
