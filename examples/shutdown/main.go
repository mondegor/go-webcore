package main

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
)

func main() {
	logger := mrcore.Log().With("shutdown")
	appHelper := mrtool.NewAppHelper(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go appHelper.GracefulShutdown(cancel)

	logger.Info("Waiting for command. To exit press CTRL+C")

	<-ctx.Done()
	logger.Info("Application stopped")
}
