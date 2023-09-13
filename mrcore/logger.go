package mrcore

import (
    "fmt"
    "strings"
)

const (
    LogErrorLevel LogLevel = iota
    LogWarnLevel
    LogInfoLevel
    LogDebugLevel
)

type (
    LogLevel uint32

    Logger interface {
        With(name string) Logger
        Error(message string, args ...any)
        Err(err error)
        Warn(message string, args ...any)
        Info(message string, args ...any)
        Debug(message string, args ...any)
    }
)

var (
	defaultLogger Logger = newLogger("[mrcore] ", LogDebugLevel)
)

// SetDefaultLogger - WARNING!!! only for main process
func SetDefaultLogger(logger Logger) {
    defaultLogger = logger
}

func DefaultLogger() Logger {
    return defaultLogger
}

func LogError(message string, args ...any) {
    if defaultLogger == nil {
        return
    }

    defaultLogger.Error(message, args...)
}

func LogErr(e error) {
    if defaultLogger == nil {
        return
    }

    defaultLogger.Err(e)
}

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
