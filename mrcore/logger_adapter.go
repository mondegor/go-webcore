package mrcore

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "time"
)

const (
    datetime = "2006/01/02 15:04:05"
)

type (
    LoggerAdapter struct {
        name    string
        level   LogLevel
        infoLog *log.Logger
        errLog *log.Logger
    }
)

// Make sure the LoggerAdapter conforms with the Logger interface
var _ Logger = (*LoggerAdapter)(nil)

// Make sure the LoggerAdapter conforms with the EventBox interface
var _ EventBox = (*LoggerAdapter)(nil)

func NewLogger(prefix string, level string) (*LoggerAdapter, error) {
    lvl, err := ParseLogLevel(level)

    if err != nil {
        return nil, err
    }

    return newLogger(prefix, lvl), nil
}

func newLogger(prefix string, level LogLevel) *LoggerAdapter {
    return &LoggerAdapter{
        level: level,
        infoLog: log.New(os.Stdout, prefix, 0),
        errLog: log.New(os.Stderr, prefix, 0),
    }
}

func (l *LoggerAdapter) With(name string) Logger {
    if l.name != "" {
        name = fmt.Sprintf("%s; %s", l.name, name)
    }

    return &LoggerAdapter{
        name: name,
        level: l.level,
        infoLog: l.infoLog,
        errLog: l.errLog,
    }
}

func (l *LoggerAdapter) Error(message string, args ...any) {
    l.logPrint(l.errLog, 3, "ERROR", message, args)
}

func (l *LoggerAdapter) Err(e error) {
    l.logPrint(l.errLog, 3, "ERROR", e.Error(), []any{})
}

func (l *LoggerAdapter) Warn(message string, args ...any) {
    if l.level >= LogWarnLevel {
        l.logPrint(l.infoLog, 3, "WARN", message, args)
    }
}

func (l *LoggerAdapter) Info(message string, args ...any) {
    if l.level >= LogInfoLevel {
        l.logPrint(l.infoLog, 0, "INFO", message, args)
    }
}

func (l *LoggerAdapter) Debug(message string, args ...any) {
    if l.level >= LogDebugLevel {
        l.logPrint(l.infoLog, 0, "DEBUG", message, args)
    }
}

func (l *LoggerAdapter) Emit(message string, args ...any) {
    l.logPrint(l.infoLog, 0, "EVENT", message, args)
}

func (l *LoggerAdapter) logPrint(logger *log.Logger, callerSkip int, prefix string, message string, args []any) {
    var buf strings.Builder

    buf.Grow(len(message) + len(l.name) + len(prefix) + 5) // 5 - separated chars
    l.formatHeader(&buf, prefix, callerSkip)
    buf.WriteString(message)

    if len(args) == 0 {
        logger.Print(buf.String())
    } else {
        logger.Printf(buf.String(), args...)
    }
}

func (l *LoggerAdapter) formatHeader(buf *strings.Builder, prefix string, callerSkip int) {
    buf.WriteString(time.Now().Format(datetime))
    buf.WriteByte(' ')

    if l.name != "" {
        buf.WriteByte('[')
        buf.WriteString(l.name)
        buf.Write([]byte{']', ' '})
    }

    buf.WriteString(prefix)
    buf.WriteByte('\t')

    if callerSkip > 0 {
        _, file, line, ok := runtime.Caller(callerSkip)

        if !ok {
            file = "???"
            line = 0
        }

        buf.WriteString(fmt.Sprintf("%s:%d\t", filepath.Base(file), line))
    }
}
