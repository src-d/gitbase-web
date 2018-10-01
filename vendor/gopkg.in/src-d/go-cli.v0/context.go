package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/src-d/go-log.v1"
)

// ContextCommander is a cancellable commander.
type ContextCommander interface {
	// ExecuteContext executes the command with a context.
	ExecuteContext(context.Context, []string) error
}

func setupContext() (context.Context, context.CancelFunc) {
	var (
		sigterm = make(chan os.Signal)
		sigint  = make(chan os.Signal)
	)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-sigterm:
			log.Infof("signal SIGTERM received, stopping...")
			cancel()
		case <-sigint:
			log.Infof("signal SIGINT received, stopping...")
			cancel()
		case <-ctx.Done():
		}

		signal.Stop(sigterm)
		signal.Stop(sigint)
	}()

	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigint, os.Interrupt)

	return ctx, cancel
}
