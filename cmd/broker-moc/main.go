package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/snyk-playground/broker-moc/internal/app"
	"github.com/snyk-playground/broker-moc/internal/command"
)

const (
	exitCodeOK        = 0
	exitCodeError     = 1
	exitCodeInterrupt = 2
)

var (
	version = "dev"
)

func main() {
	code := mainRun()
	os.Exit(code)
}

func mainRun() int {
	ctx, cancel := handleInterrupt(context.Background())
	defer cancel()

	cfg, err := app.NewConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return exitCodeError
	}

	bma, err := app.New(cfg, version)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return exitCodeError
	}
	rootCmd := command.NewRootCmd(bma)
	if err = rootCmd.ExecuteContext(ctx); err != nil {
		return exitCodeError
	}

	return exitCodeOK
}

// try to exit gracefully when the interrupt signal is sent (CTRL+C).
func handleInterrupt(ctx context.Context) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		select {
		case <-signalChan:
			// first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan
		// second signal, hard exit
		os.Exit(exitCodeInterrupt)
	}()

	return ctx, func() {
		signal.Stop(signalChan)
		cancel()
	}
}
