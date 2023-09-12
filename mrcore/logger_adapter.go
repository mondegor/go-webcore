package mrcore

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
    loggerAdapter struct {
        name string
        level LogLevel
        infoLog *log.Logger
        errLog *log.Logger
    }
)

func NewLogger(prefix string, level string) (Logger, error) {
    lvl, err := ParseLogLevel(level)

    if err != nil {
        return nil, err
    }

    return newLogger(prefix, lvl), nil
}

func newLogger(prefix string, level LogLevel) Logger {
    return &loggerAdapter {
        level: level,
        infoLog: log.New(os.Stdout, prefix, 0),
        errLog: log.New(os.Stderr, prefix, 0),
    }
}

func (l *loggerAdapter) With(name string) Logger {
    if l.name != "" {
        name = fmt.Sprintf("%s:%s", l.name, name)
    }

    return &loggerAdapter {
        name: name,
        level: l.level,
        infoLog: l.infoLog,
        errLog: l.errLog,
    }
}

func (l *loggerAdapter) Error(message string, args ...any) {
    l.logPrint(l.errLog, 3, "ERROR", message, args)
}

func (l *loggerAdapter) Err(e error) {
    l.logPrint(l.errLog, 3, "ERROR", e.Error(), []any{})
}

func (l *loggerAdapter) Warn(message string, args ...any) {
    if l.level >= LogWarnLevel {
        l.logPrint(l.infoLog, 3, "WARN", message, args)
    }
}

func (l *loggerAdapter) Info(message string, args ...any) {
    if l.level >= LogInfoLevel {
        l.logPrint(l.infoLog, 0, "INFO", message, args)
    }
}

func (l *loggerAdapter) Debug(message string, args ...any) {
    if l.level >= LogDebugLevel {
        l.logPrint(l.infoLog, 0, "DEBUG", message, args)
    }
}

func (l *loggerAdapter) Emit(message string, args ...any) {
    l.logPrint(l.infoLog, 0, "EVENT", message, args)
}

func (l *loggerAdapter) logPrint(logger *log.Logger, callerSkip int, prefix string, message string, args []any) {
    var buf []byte

    l.formatHeader(&buf, prefix, callerSkip)

    buf = append(buf, message...)

    if len(args) == 0 {
        logger.Print(string(buf))
    } else {
        logger.Printf(string(buf), args...)
    }
}

func (l *loggerAdapter) formatHeader(buf *[]byte, prefix string, callerSkip int) {
    *buf = append(*buf, time.Now().Format(datetime)...)
    *buf = append(*buf, ' ')

    if l.name != "" {
        *buf = append(*buf, '[')
        *buf = append(*buf, l.name...)
        *buf = append(*buf, ']', ' ')
    }

    *buf = append(*buf, prefix...)
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
