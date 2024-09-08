package mrlog

import (
	"fmt"
	"log"
)

type (
	eventAdapter struct {
		logger *log.Logger
		buf    []byte
		done   func(msg string)
	}
)

// Make sure the Image conforms with the mrlog.LoggerEvent interface.
var _ LoggerEvent = (*eventAdapter)(nil)

// Err - comment method.
func (e *eventAdapter) Err(err error) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(c.buf, "err"), err.Error())

	return &c
}

// Str - comment method.
func (e *eventAdapter) Str(key, value string) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(e.buf, key), value)

	return &c
}

// Bytes - comment method.
func (e *eventAdapter) Bytes(key string, value []byte) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(e.buf, key), value)

	return &c
}

// Int - comment method.
func (e *eventAdapter) Int(key string, value int) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(e.buf, key), value)

	return &c
}

// Any - comment method.
func (e *eventAdapter) Any(key string, value any) LoggerEvent {
	if e == nil {
		return e
	}

	c := *e
	c.buf = appendValue(appendKey(e.buf, key), value)

	return &c
}

// Msg - comment method.
func (e *eventAdapter) Msg(message string) {
	if e == nil {
		return
	}

	e.write(message + string(e.buf))
}

// Msgf - comment method.
func (e *eventAdapter) Msgf(format string, args ...any) {
	if e == nil {
		return
	}

	e.write(fmt.Sprintf(format+string(e.buf), args...))
}

// MsgFunc - comment method.
func (e *eventAdapter) MsgFunc(createMsg func() string) {
	if e == nil {
		return
	}

	e.write(createMsg() + string(e.buf))
}

// Send - comment method.
func (e *eventAdapter) Send() {
	if e == nil {
		return
	}

	e.write("[empty]" + string(e.buf))
}

func (e *eventAdapter) write(message string) {
	e.logger.Print(message)

	if e.done != nil {
		e.done(message)
	}
}
