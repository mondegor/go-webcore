package main

import (
    "context"

    "github.com/mondegor/go-webcore/mrcore"
)

func main() {
    logger := mrcore.DefaultLogger().With("shutdown")
    appHelper := mrcore.NewAppHelper(logger)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go appHelper.GracefulShutdown(cancel)

    logger.Info("Waiting for command. To exit press CTRL+C")

    <-ctx.Done()
    logger.Info("Application stopped")
}
