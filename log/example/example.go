package main

import (
	"context"
	e "errors"
	"log/slog"
	"os"

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
	defer os.Exit(0)

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

	defer log.Info(ctx, "Deferred")

	attrs := []slog.Attr{
		log.String("attr1", "value1"),
		log.String("attr2", "value2"),
	}

	log.Debug(ctx, "Debug", attrs...)
	log.Info(ctx, "Info", attrs...)
	log.Warn(ctx, "Warn", attrs...)
	log.Errorw(ctx, "Error", attrs...)
	log.Criticalw(ctx, "Critical", attrs...)

	log.Error(ctx, errDefault.WithStack(), attrs...)
	log.Critical(ctx, errCritical.WithStack(), attrs...)
	log.Fatal(ctx, errFatal.WithStack(), attrs...)
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
			Retryable(true)

	errFatal = errors.New("fatal error").
			WithKind(errors.KindUnexpected).
			WithCode("ERR_FATAL").
			WithRootError(e.New("root fatal error"))
)
