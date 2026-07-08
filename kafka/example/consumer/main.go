package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/lcnascimento/go-kit/kafka"
	"github.com/lcnascimento/go-kit/kafka/example"
	"github.com/lcnascimento/go-kit/o11y"
)

func main() {
	ctx := context.Background()
	o11y.MustStart(ctx)
	defer o11y.Shutdown(ctx)

	subscriber := kafka.NewSubscriber[example.Example]("example")
	defer subscriber.Stop(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, os.Kill)
		<-sigCh
		cancel()
	}()

	subscriber.Run(ctx, func(ctx context.Context, e example.Example) error {
		slog.InfoContext(ctx, e.Message)
		return nil
	})
}
