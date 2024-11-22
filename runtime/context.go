package runtime

import (
	"context"
	"os"
	"os/signal"
)

// ContextWithOSSignalCancellation builds a new Context that cancels itself on os.Interrupt signals.
func ContextWithOSSignalCancellation() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		cancel()
	}()

	return ctx
}
