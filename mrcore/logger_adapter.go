package mrcore

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
)

const (
	datetime = "2006/01/02 15:04:05"
)

type (
	LoggerAdapter struct {
		name              string
		level             LogLevel
		caller            *mrerr.Caller
		callerSkip        int
		callerEnabledFunc func(err error) bool
		infoLog           *log.Logger
		errLog            *log.Logger
	}

	LoggerOptions struct {
		Prefix            string
		Level             string
		CallerOptions     []mrerr.CallerOption
		CallerEnabledFunc func(err error) bool
	}
)

// Make sure the LoggerAdapter conforms with the Logger interface
var _ Logger = (*LoggerAdapter)(nil)

// Make sure the LoggerAdapter conforms with the EventBox interface
var _ EventBox = (*LoggerAdapter)(nil)

func NewLogger(opt LoggerOptions) (*LoggerAdapter, error) {
	level, err := ParseLogLevel(opt.Level)

	if err != nil {
		return nil, err
	}

	return newLogger(opt, level), nil
}

func newLogger(opt LoggerOptions, level LogLevel) *LoggerAdapter {
	l := LoggerAdapter{
		level:             level,
		caller:            mrerr.NewCaller(opt.CallerOptions...),
		callerEnabledFunc: opt.CallerEnabledFunc,
		callerSkip:        3, // skip: .., logPrint, formatHeader
		infoLog:           log.New(os.Stdout, opt.Prefix, 0),
		errLog:            log.New(os.Stderr, opt.Prefix, 0),
	}

	if l.callerEnabledFunc == nil {
		l.callerEnabledFunc = func(err error) bool {
			return true
		}
	}

	return &l
}

func (l *LoggerAdapter) With(name string) Logger {
	if name == "" {
		return l
	}

	if l.name != "" {
		name = l.name + ";" + name
	}

	c := *l
	c.name = name

	return &c
}

func (l *LoggerAdapter) Caller(skip int) Logger {
	if skip == 0 {
		return l
	}

	c := *l
	c.callerSkip += skip

	return &c
}

func (l *LoggerAdapter) Level() LogLevel {
	return l.level
}

func (l *LoggerAdapter) Error(message string, args ...any) {
	l.logPrint(l.errLog, "ERROR", message, args, true)
}

func (l *LoggerAdapter) Err(e error) {
	if e == nil {
		return
	}

	l.logPrint(l.errLog, "ERROR", e.Error(), []any{}, l.callerEnabledFunc(e))
}

func (l *LoggerAdapter) Warning(message string, args ...any) {
	if l.level >= LogWarnLevel {
		l.logPrint(l.errLog, "WARN", message, args, true)
	}
}

func (l *LoggerAdapter) Warn(e error) {
	if l.level >= LogWarnLevel && e != nil {
		l.logPrint(l.errLog, "WARN", e.Error(), []any{}, l.callerEnabledFunc(e))
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

func (l *LoggerAdapter) logPrint(logger *log.Logger, prefix, message string, args []any, showCallStack bool) {
	var buf strings.Builder

	buf.Grow(len(message) + len(l.name) + len(prefix) + 24) // 24 = (19 - datetime, 5 - separated chars)
	l.formatHeader(&buf, prefix, showCallStack)
	buf.WriteString(message)

	if len(args) == 0 {
		logger.Print(buf.String())
	} else {
		logger.Printf(buf.String(), args...)
	}
}

func (l *LoggerAdapter) formatHeader(buf *strings.Builder, prefix string, showCallStack bool) {
	buf.WriteString(time.Now().UTC().Format(datetime))
	buf.WriteByte(' ')

	if l.name != "" {
		buf.WriteByte('[')
		buf.WriteString(l.name)
		buf.Write([]byte{']', ' '})
	}

	buf.WriteString(prefix)
	buf.WriteByte('\t')

	if showCallStack {
		cs := l.caller.CallStack(l.callerSkip)

		if len(cs) > 0 {
			for i := range cs {
				if i > 0 {
					buf.Write([]byte{' ', '<', '-', ' '})
				}

				buf.WriteString(cs[i].File)
				buf.WriteByte(':')
				buf.WriteString(strconv.Itoa(cs[i].Line))
			}

			buf.WriteByte('\t')
		}
	}
}
