package main

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/mrlogbase"
	"github.com/mondegor/go-webcore/mrserver"
)

func main() {
	logger := mrlogbase.New(mrlog.DebugLevel).With().Str("example", "shutdown").Logger()

	ctx, cancel := context.WithCancel(logger.WithContext(context.Background()))
	defer cancel()

	exec, intr := mrserver.PrepareAppToStart(ctx)
	defer intr(nil)

	if err := exec(); err != nil {
		logger.Info().Msg("Application stopped with error")
	} else {
		logger.Info().Msg("Application stopped")
	}
}
