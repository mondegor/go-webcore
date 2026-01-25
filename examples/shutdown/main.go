package main

import (
	"context"
	"os"
	"sync/atomic"
	"time"

	"github.com/mondegor/go-sysmess/mrlog/slog"

	"github.com/mondegor/go-webcore/mrworker/process/signal"
)

var isTimeout atomic.Bool

func main() {
	logger, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := signal.NewInterceptor(logger)

	defer func() {
		if !isTimeout.Swap(true) {
			_ = app.Shutdown(ctx)
			logger.Info(ctx, "The application has been stopped")
		}
	}()

	timer := time.After(3 * time.Second)

	go func() {
		select {
		case <-ctx.Done():
		case <-timer:
			cancel()
		case <-timer:
			if !isTimeout.Swap(true) {
				_ = app.Shutdown(ctx)
				logger.Error(ctx, "The application has been interrupted by timeout")
			}
		}
	}()

	logger.Info(ctx, "The application started, waiting for requests. To exit press CTRL+C")

	if err := app.Start(ctx, nil); err != nil {
		logger.Error(ctx, "Application has been interrupted with error", "error", err)
	}

	// задержка нужна, чтобы залогировались данные из горутины
	if isTimeout.Load() {
		time.Sleep(time.Second)
	}
}
