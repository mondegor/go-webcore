package main

import (
    "flag"
    "fmt"

    "github.com/mondegor/go-core/mrlog"
)

var level string

func init() {
    flag.StringVar(&level, "level", "info", "Logging level")
}

func main() {
    flag.Parse()

    logger, err := mrlog.New("[example]", level, true)

    if err != nil {
        fmt.Println(err)
        return
    }

    logger.Info("Test message")
}
