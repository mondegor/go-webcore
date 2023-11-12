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
        Caller(skip int) Logger
        DisableFileLine() Logger
        Error(message string, args ...any)
        Err(err error)
        Warning(message string, args ...any)
        Warn(err error)
        Info(message string, args ...any)
        Debug(message string, args ...any)
    }
)

var (
    defaultLogger Logger = newLogger("[mrcore] ", LogDebugLevel)
)

// SetDefaultLogger - WARNING: use only when starting the main process
func SetDefaultLogger(logger Logger) {
    defaultLogger = logger
}

func Log() Logger {
    return defaultLogger
}

func LogError(message string, args ...any) {
    if defaultLogger == nil {
        return
    }

    defaultLogger.Caller(1).Error(message, args...)
}

func LogErr(e error) {
    if defaultLogger == nil {
        return
    }

    defaultLogger.Caller(1).Err(e)
}

func LogWarning(message string, args ...any) {
    if defaultLogger == nil {
        return
    }

    defaultLogger.Caller(1).Warning(message, args...)
}

func LogWarn(e error) {
    if defaultLogger == nil {
        return
    }

    defaultLogger.Caller(1).Warn(e)
}

func LogInfo(message string, args ...any) {
    if defaultLogger == nil {
        return
    }

    defaultLogger.Info(message, args...)
}

func LogDebug(message string, args ...any) {
    if defaultLogger == nil {
        return
    }

    defaultLogger.Debug(message, args...)
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
