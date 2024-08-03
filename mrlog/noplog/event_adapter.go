package noplog

import (
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	eventAdapter struct {
		done func()
	}
)

// Make sure the Image conforms with the mrlog.LoggerEvent interface.
var _ mrlog.LoggerEvent = (*eventAdapter)(nil)

// Err - comment method.
func (e *eventAdapter) Err(_ error) mrlog.LoggerEvent {
	return e
}

// Str - comment method.
func (e *eventAdapter) Str(_, _ string) mrlog.LoggerEvent {
	return e
}

// Bytes - comment method.
func (e *eventAdapter) Bytes(_ string, _ []byte) mrlog.LoggerEvent {
	return e
}

// Int - comment method.
func (e *eventAdapter) Int(_ string, _ int) mrlog.LoggerEvent {
	return e
}

// Any - comment method.
func (e *eventAdapter) Any(_ string, _ any) mrlog.LoggerEvent {
	return e
}

// Msg - comment method.
func (e *eventAdapter) Msg(_ string) {
	if e == nil {
		return
	}

	e.execDone()
}

// Msgf - comment method.
func (e *eventAdapter) Msgf(_ string, _ ...any) {
	if e == nil {
		return
	}

	e.execDone()
}

// MsgFunc - comment method.
func (e *eventAdapter) MsgFunc(_ func() string) {
	if e == nil {
		return
	}

	e.execDone()
}

// Send - comment method.
func (e *eventAdapter) Send() {
	if e == nil {
		return
	}

	e.execDone()
}

func (e *eventAdapter) execDone() {
	if e.done != nil {
		e.done()
	}
}
