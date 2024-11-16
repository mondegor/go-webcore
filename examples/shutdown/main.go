package main

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker/process/signal"
)

func main() {
	logger := mrlog.Default().With().Str("example", "shutdown").Logger()
	ctx, app := signal.NewInterception(logger.WithContext(context.Background()))

	defer func() {
		_ = app.Shutdown(ctx)
		mrlog.Ctx(ctx).Info().Msg("The application has been stopped")
	}()

	mrlog.Ctx(ctx).Info().Msg("The application started, waiting for requests. To exit press CTRL+C")

	if err := app.Start(ctx, nil); err != nil {
		logger.Info().Msg("Application stopped with error")
	} else {
		logger.Info().Msg("Application stopped")
	}
}
