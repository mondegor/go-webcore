package main

import (
	"time"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/mrzerolog"
)

func main() {
	logger := mrzerolog.New(
		mrlog.Options{
			Level:           mrlog.InfoLevel,
			JsonFormat:      true,
			TimestampFormat: time.RFC3339Nano,
		},
	)

	mrlog.SetDefault(logger)

	logger.Info().Msg("Logger info message - OK!")
	logger.Debug().Msg("Logger debug message skipped")
}
