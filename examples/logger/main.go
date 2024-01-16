package main

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/go-webcore/mrcore"
)

func main() {
	logger, err := mrcore.NewLogger(
		mrcore.LoggerOptions{
			Prefix: "[my-logger] ",
			Level:  "info",
			CallerOptions: []mrerr.CallerOption{
				mrerr.CallerDeep(4),
				mrerr.CallerUseShortPath(true),
				mrerr.CallerRootPath("/"),
			},
		},
	)

	if err != nil {
		fmt.Println("create logger error")
		return
	}

	mrcore.SetDefaultLogger(logger)

	mrcore.LogInfo("Logger info message - OK!")
	mrcore.LogDebug("Logger debug message skipped")
}
