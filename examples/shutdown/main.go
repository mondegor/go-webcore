package main

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrrun"
)

func main() {
	logger := mrlog.Default().With().Str("example", "shutdown").Logger()

	_, executer := mrrun.MakeSignalHandler(logger.WithContext(context.Background()))
	defer executer.Interrupt(nil)

	if err := executer.Execute(); err != nil {
		logger.Info().Msg("Application stopped with error")
	} else {
		logger.Info().Msg("Application stopped")
	}
}
