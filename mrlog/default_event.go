package mrlog

import (
	"fmt"
	"log"
)

type (
	defaultEvent struct {
		logger *log.Logger
		buf    []byte
		done   func(msg string)
	}
)

func (e *defaultEvent) Caller(skip ...int) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	return &c
}

func (e *defaultEvent) Err(err error) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(c.buf, "err"), err.Error())
	return &c
}

func (e *defaultEvent) Str(key, value string) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(e.buf, key), value)
	return &c
}

func (e *defaultEvent) Bytes(key string, value []byte) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(e.buf, key), value)
	return &c
}

func (e *defaultEvent) Int(key string, value int) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(e.buf, key), value)
	return &c
}

func (e *defaultEvent) Any(key string, value any) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(e.buf, key), value)
	return &c
}

func (e *defaultEvent) Msg(message string) {
	if e == nil {
		return
	}

	e.write(defaultLoggerPrefix + message + string(e.buf))
}

func (e *defaultEvent) Msgf(format string, args ...any) {
	if e == nil {
		return
	}

	e.write(fmt.Sprintf(defaultLoggerPrefix+format+string(e.buf), args...))
}

func (e *defaultEvent) MsgFunc(createMsg func() string) {
	if e == nil {
		return
	}

	e.write(defaultLoggerPrefix + createMsg() + string(e.buf))
}

func (e *defaultEvent) Send() {
	if e == nil {
		return
	}

	e.write(defaultLoggerPrefix[:len(defaultLoggerPrefix)-1] + string(e.buf))
}

func (e *defaultEvent) write(message string) {
	if e.done != nil {
		defer e.done(message)
	}

	e.logger.Print(message)
}
