package mrserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mondegor/go-webcore/mrlog"
)

const (
	signalChanLen = 10
)

// PrepareAppToStart - comment func.
func PrepareAppToStart(ctx context.Context) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)
	signalStop := make(chan os.Signal, signalChanLen)

	signal.Notify(
		signalStop,
		syscall.SIGABRT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		os.Interrupt,
		syscall.SIGTERM,
	)

	return func() error {
			logger := mrlog.Ctx(ctx)
			logger.Info().Msg("Waiting for requests. To exit press CTRL+C")

			select {
			case signalApp := <-signalStop:
				logger.Info().Msgf("Shutting down the application by signal: " + signalApp.String())

				return nil
			case <-ctx.Done():
				logger.Info().Msgf("Shutting down the application by a child process")

				return ctx.Err()
			}
		}, func(_ error) {
			cancel()
		}
}
