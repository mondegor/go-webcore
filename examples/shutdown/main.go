package main

import (
	"context"
	"math/rand/v2"
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

	timer1 := time.After(time.Duration(1000+rand.IntN(2001)) * time.Millisecond)
	timer2 := time.After(time.Duration(1000+rand.IntN(2001)) * time.Millisecond)

	go func() {
		select {
		case <-ctx.Done():
		case <-timer1:
			cancel()
		case <-timer2:
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
