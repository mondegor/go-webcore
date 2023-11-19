package mrcore

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	datetime = "2006/01/02 15:04:05"
)

type (
	LoggerAdapter struct {
		name            string
		level           LogLevel
		callerSkip      int
		enabledFileLine bool
		infoLog         *log.Logger
		errLog          *log.Logger
	}
)

// Make sure the LoggerAdapter conforms with the Logger interface
var _ Logger = (*LoggerAdapter)(nil)

// Make sure the LoggerAdapter conforms with the EventBox interface
var _ EventBox = (*LoggerAdapter)(nil)

func NewLogger(prefix, level string) (*LoggerAdapter, error) {
	lvl, err := ParseLogLevel(level)

	if err != nil {
		return nil, err
	}

	return newLogger(prefix, lvl), nil
}

func newLogger(prefix string, level LogLevel) *LoggerAdapter {
	return &LoggerAdapter{
		level:           level,
		callerSkip:      4, // to parent function
		enabledFileLine: true,
		infoLog:         log.New(os.Stdout, prefix, 0),
		errLog:          log.New(os.Stderr, prefix, 0),
	}
}

func (l LoggerAdapter) With(name string) Logger {
	if l.name != "" {
		name = l.name + "; " + name
	}

	l.name = name
	return &l
}

func (l LoggerAdapter) Caller(skip int) Logger {
	l.callerSkip += skip
	return &l
}

func (l LoggerAdapter) DisableFileLine() Logger {
	l.enabledFileLine = false
	return &l
}

func (l *LoggerAdapter) Error(message string, args ...any) {
	l.logPrint(l.errLog, "ERROR", message, args, true)
}

func (l *LoggerAdapter) Err(e error) {
	if e == nil {
		return
	}

	l.logPrint(l.errLog, "ERROR", e.Error(), []any{}, true)
}

func (l *LoggerAdapter) Warning(message string, args ...any) {
	if l.level >= LogWarnLevel {
		l.logPrint(l.errLog, "WARN", message, args, true)
	}
}

func (l *LoggerAdapter) Warn(e error) {
	if l.level >= LogWarnLevel && e != nil {
		l.logPrint(l.errLog, "WARN", e.Error(), []any{}, true)
	}
}

func (l *LoggerAdapter) Info(message string, args ...any) {
	if l.level >= LogInfoLevel {
		l.logPrint(l.infoLog, "INFO", message, args, false)
	}
}

func (l *LoggerAdapter) Debug(message string, args ...any) {
	if l.level >= LogDebugLevel {
		l.logPrint(l.infoLog, "DEBUG", message, args, false)
	}
}

func (l *LoggerAdapter) Emit(message string, args ...any) {
	l.logPrint(l.infoLog, "EVENT", message, args, false)
}

func (l *LoggerAdapter) logPrint(logger *log.Logger, prefix, message string, args []any, showFileLine bool) {
	var buf strings.Builder

	buf.Grow(len(message) + len(l.name) + len(prefix) + 24) // 24 = (19 - datetime, 5 - separated chars)
	l.formatHeader(&buf, prefix, showFileLine)
	buf.WriteString(message)

	if len(args) == 0 {
		logger.Print(buf.String())
	} else {
		logger.Printf(buf.String(), args...)
	}
}

func (l *LoggerAdapter) formatHeader(buf *strings.Builder, prefix string, showFileLine bool) {
	buf.WriteString(time.Now().Format(datetime))
	buf.WriteByte(' ')

	if l.name != "" {
		buf.WriteByte('[')
		buf.WriteString(l.name)
		buf.Write([]byte{']', ' '})
	}

	buf.WriteString(prefix)
	buf.WriteByte('\t')

	if l.enabledFileLine && showFileLine {
		_, file, line, ok := runtime.Caller(l.callerSkip)

		if !ok {
			file = "???"
			line = 0
		}

		buf.WriteString(file)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
		buf.WriteByte('\t')
	}
}
