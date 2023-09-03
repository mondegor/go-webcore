package mrlog

import (
    "fmt"
    "strings"
)

type LogLevel uint32

const (
    LogErrorLevel LogLevel = iota
    LogWarnLevel
    LogInfoLevel
    LogDebugLevel
)

func ParseLogLevel(level string) (LogLevel, error) {
    switch strings.ToLower(level) {
    case "error":
        return LogErrorLevel, nil

    case "warn", "warning":
        return LogWarnLevel, nil

    case "info":
        return LogInfoLevel, nil

    case "debug":
        return LogDebugLevel, nil
    }

    return LogErrorLevel, fmt.Errorf("log level '%s' is unknown", level)
}
