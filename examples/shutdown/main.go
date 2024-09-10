package main

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrapp"
)

func main() {
	logger := mrlog.Default().With().Str("example", "shutdown").Logger()

	ctx, cancel := context.WithCancel(logger.WithContext(context.Background()))
	defer cancel()

	exec, intr := mrapp.PrepareToStart(ctx)
	defer intr(nil)

	if err := exec(); err != nil {
		logger.Info().Msg("Application stopped with error")
	} else {
		logger.Info().Msg("Application stopped")
	}
}
