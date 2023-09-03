package main

import (
    "context"
    "fmt"

    "github.com/mondegor/go-core/mrcore"
    "github.com/mondegor/go-core/mrlog"
)

func main() {
    logger, err := mrlog.New("[shutdown]", "info", true)

    if err != nil {
        fmt.Println(err)
        return
    }

    appHelper := mrcore.NewAppHelper(logger)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go appHelper.GracefulShutdown(cancel)

    logger.Info("Waiting for command. To exit press CTRL+C")

    <-ctx.Done()
    logger.Info("Application stopped")
}
