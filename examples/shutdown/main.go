package main

import (
	"context"
	"os"

	"github.com/mondegor/go-sysmess/mrlog/slog"

	"github.com/mondegor/go-webcore/mrworker/process/signal"
)

func main() {
	logger, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))

	ctx, app := signal.NewInterceptor(context.Background(), logger)

	defer func() {
		_ = app.Shutdown(ctx)
		logger.Info(ctx, "The application has been stopped")
	}()

	logger.Info(ctx, "The application started, waiting for requests. To exit press CTRL+C")

	if err := app.Start(ctx, nil); err != nil {
		logger.Info(ctx, "Application stopped with error")
	} else {
		logger.Info(ctx, "Application stopped")
	}
}
