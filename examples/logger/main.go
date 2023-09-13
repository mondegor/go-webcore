package main

import (
    "fmt"

    "github.com/mondegor/go-webcore/mrcore"
)

func main() {
    logger, err := mrcore.NewLogger("[my-logger]", "info")

    if err != nil {
        fmt.Println("create logger error")
        return
    }

    mrcore.SetDefaultLogger(logger)

    mrcore.DefaultLogger().Info("Logger info message - OK!")
    mrcore.DefaultLogger().Debug("Logger debug message skipped")
}
