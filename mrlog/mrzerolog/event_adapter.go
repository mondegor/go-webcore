package mrzerolog

import (
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/rs/zerolog"
)

// go get -u github.com/rs/zerolog

type (
	EventAdapter struct {
		ze                  *zerolog.Event
		isAutoCallerAllowed bool
		isAutoCallerOnFunc  func(err error) bool
	}
)

func (e *EventAdapter) Caller(skip ...int) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	if len(skip) > 0 {
		skip[0]++
	}

	ev := e.newEventAdapter(e.ze.Caller(skip...))
	ev.isAutoCallerAllowed = false

	return ev
}

func (e *EventAdapter) CallerSkipFrame(count int) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return e.newEventAdapter(e.ze.CallerSkipFrame(count + 1))
}

func (e *EventAdapter) Err(err error) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	ev := e.newEventAdapter(e.ze.Err(err))

	if ev.isAutoCallerAllowed {
		if e.isAutoCallerOnFunc == nil || !e.isAutoCallerOnFunc(err) {
			ev.isAutoCallerAllowed = false
		}
	}

	return ev
}

func (e *EventAdapter) Str(key, value string) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return e.newEventAdapter(e.ze.Str(key, value))
}

func (e *EventAdapter) Bytes(key string, value []byte) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return e.newEventAdapter(e.ze.Bytes(key, value))
}

func (e *EventAdapter) Int(key string, value int) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return e.newEventAdapter(e.ze.Int(key, value))
}

func (e *EventAdapter) Any(key string, value any) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return e.newEventAdapter(e.ze.Any(key, value))
}

func (e *EventAdapter) Msg(message string) {
	if e == nil {
		return
	}

	e.prepareEvent().Msg(message)
}

func (e *EventAdapter) Msgf(format string, args ...any) {
	if e == nil {
		return
	}

	e.prepareEvent().Msgf(format, args...)
}

func (e *EventAdapter) MsgFunc(createMsg func() string) {
	if e == nil {
		return
	}

	e.prepareEvent().MsgFunc(createMsg)
}

func (e *EventAdapter) Send() {
	if e == nil {
		return
	}

	e.prepareEvent().Send()
}

func (e *EventAdapter) newEventAdapter(ze *zerolog.Event) *EventAdapter {
	c := *e
	c.ze = ze

	return &c
}

func (e *EventAdapter) prepareEvent() *zerolog.Event {
	if e.isAutoCallerAllowed {
		return e.ze.Caller(2)
	}

	return e.ze
}
