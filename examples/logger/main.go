package main

import (
	"fmt"

	"github.com/mondegor/go-webcore/mrcore"
)

func main() {
	logger, err := mrcore.NewLogger(
		mrcore.LoggerOptions{
			Prefix: "[my-logger] ",
			Level:  "info",
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
