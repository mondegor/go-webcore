package zerolog

import (
	"github.com/rs/zerolog"

	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// eventAdapter - comment struct.
	eventAdapter struct {
		ze               *zerolog.Event
		prepareErrorFunc func(err error) error
	}
)

// Err - comment method.
func (e *eventAdapter) Err(err error) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	err = e.prepareErrorFunc(err)

	return e.newEventAdapter(e.ze.Err(err))
}

// Str - comment method.
func (e *eventAdapter) Str(key, value string) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return e.newEventAdapter(e.ze.Str(key, value))
}

// Bytes - comment method.
func (e *eventAdapter) Bytes(key string, value []byte) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return e.newEventAdapter(e.ze.Bytes(key, value))
}

// Int - comment method.
func (e *eventAdapter) Int(key string, value int) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return e.newEventAdapter(e.ze.Int(key, value))
}

// Any - comment method.
func (e *eventAdapter) Any(key string, value any) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return e.newEventAdapter(e.ze.Any(key, value))
}

// Msg - comment method.
func (e *eventAdapter) Msg(message string) {
	if e == nil {
		return
	}

	e.ze.Msg(message)
}

// Msgf - comment method.
func (e *eventAdapter) Msgf(format string, args ...any) {
	if e == nil {
		return
	}

	e.ze.Msgf(format, args...)
}

// MsgFunc - comment method.
func (e *eventAdapter) MsgFunc(createMsg func() string) {
	if e == nil {
		return
	}

	e.ze.MsgFunc(createMsg)
}

// Send - comment method.
func (e *eventAdapter) Send() {
	if e == nil {
		return
	}

	e.ze.Send()
}

func (e *eventAdapter) newEventAdapter(ze *zerolog.Event) *eventAdapter {
	c := *e
	c.ze = ze

	return &c
}
