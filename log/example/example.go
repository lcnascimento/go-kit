package main

import (
	"context"
	"fmt"

	"github.com/lcnascimento/go-kit/log"
	"github.com/lcnascimento/go-kit/log/format"
	"github.com/lcnascimento/go-kit/propagation"
)

func main() {
	ctx := context.Background()

	contextKeys := propagation.ContextKeySet{
		propagation.ContextKey("foo"): true,
	}

	logger := log.NewLogger(
		log.WithLevel("DEBUG"),
		log.WithContextKeySet(contextKeys),
	)

	ctx = context.WithValue(ctx, propagation.ContextKey("foo"), "bar")

	attrs := format.AttributeSet{
		"attr1": "value1",
		"attr2": "value2",
	}

	logger.Debug(ctx, "debug message", attrs)
	logger.Info(ctx, "info message", attrs)
	logger.Warning(ctx, "warning message", attrs)
	logger.Error(ctx, fmt.Errorf("error message"), attrs)
	logger.Critical(ctx, fmt.Errorf("critical message"), attrs)
}
