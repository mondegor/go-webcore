package mrlog

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "time"
)

const (
    datetime = "2006/01/02 15:04:05"
)

type (
    Logger interface {
        With(name string) Logger
        Error(message any, args ...any)
        Warn(message string, args ...any)
        Info(message string, args ...any)
        Debug(message any, args ...any)
    }

	logger struct {
        name string
        level LogLevel
        color bool
        infoLog *log.Logger
        errLog *log.Logger
    }
)

var std = newLogger("[default] ", LogInfoLevel, false)

func Default() Logger {
    return std
}

func Error(message any, args ...any) {
    std.logPrint(std.errLog, 3, "ERROR", message, args)
}

func New(prefix string, level string, color bool) (Logger, error) {
    lvl, err := ParseLogLevel(level)

    if err != nil {
        return nil, err
    }

    return newLogger(prefix + " ", lvl, color), nil
}

func newLogger(prefix string, level LogLevel, color bool) *logger {
    return &logger {
        level: level,
        color: color,
        infoLog: log.New(os.Stdout, prefix, 0),
        errLog: log.New(os.Stderr, prefix, 0),
    }
}

func (l *logger) With(name string) Logger {
    if l.name != "" {
        name = l.name + ":" + name
    }

    return &logger {
        name: name,
        level: l.level,
        color: l.color,
        infoLog: l.infoLog,
        errLog: l.errLog,
    }
}

func (l *logger) Error(message any, args ...any) {
    l.logPrint(l.errLog, 3, "ERROR", message, args)
}

func (l *logger) Warn(message string, args ...any) {
    if l.level >= LogWarnLevel {
        l.logPrint(l.infoLog, 3, "WARN", message, args)
    }
}

func (l *logger) Info(message string, args ...any) {
    if l.level >= LogInfoLevel {
        l.logPrint(l.infoLog, 0, "INFO", message, args)
    }
}

func (l *logger) Debug(message any, args ...any) {
    if l.level >= LogDebugLevel {
        l.logPrint(l.infoLog, 0, "DEBUG", message, args)
    }
}

func (l *logger) Emit(message string, args ...any) {
    l.logPrint(l.infoLog, 0, "EVENT", message, args)
}

func (l *logger) logPrint(logger *log.Logger, callerSkip int, prefix string, message any, args []any) {
    var buf []byte

    l.formatHeader(&buf, prefix, callerSkip)
    l.formatMessage(&buf, message)

    if len(args) == 0 {
        logger.Print(string(buf))
    } else {
        logger.Printf(string(buf), args...)
    }
}

func (l *logger) formatHeader(buf *[]byte, prefix string, callerSkip int) {
    *buf = append(*buf, time.Now().Format(datetime)...)
    *buf = append(*buf, ' ')
    *buf = append(*buf, prefix...)

    if l.name != "" {
        *buf = append(*buf, ' ', '[')
        *buf = append(*buf, l.name...)
        *buf = append(*buf, ']', ' ')
    }

    *buf = append(*buf, '\t')

    if callerSkip > 0 {
        _, file, line, ok := runtime.Caller(callerSkip)

        if !ok {
            file = "???"
            line = 0
        }

        *buf = append(*buf, fmt.Sprintf("%s:%d\t", filepath.Base(file), line)...)
    }
}

func (l *logger) formatMessage(buf *[]byte, message any) {
    switch msg := message.(type) {
    case error:
        *buf = append(*buf, msg.Error()...)

    case string:
        *buf = append(*buf, msg...)

    default:
        *buf = append(*buf, fmt.Sprintf("Message %v has unknown type %v", message, msg)...)
    }
}
