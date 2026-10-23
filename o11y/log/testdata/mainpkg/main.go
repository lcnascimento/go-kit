package main

import (
	"context"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y"
	"github.com/lcnascimento/go-kit/o11y/log"
)

func main() {
	ctx := context.Background()
	o11y.MustStart(ctx)

	fail(ctx)

	if err := o11y.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func fail(ctx context.Context) {
	log.MustNewLogger("main").Error(ctx, errors.New("boom").WithKind(errors.KindInternal))
}
