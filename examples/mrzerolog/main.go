package main

import (
	"errors"
	"log"
	"os"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/zerolog/factory"
)

func main() {
	logger, err := factory.NewZeroLogAdapter(
		factory.Options{
			Stdout:          os.Stdout,
			Level:           mrlog.InfoLevel.String(),
			JsonFormat:      false,
			TimestampFormat: "RFC3339Nano",
			ConsoleColor:    true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info().Msg("Logger info message - OK!")
	logger.Debug().Msg("Logger debug message skipped")
	logger.Error().Err(errors.New("my error")).Msg("Error with auto caller")
}
