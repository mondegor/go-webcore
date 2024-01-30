package mrzerolog

import (
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/rs/zerolog"
)

// go get -u github.com/rs/zerolog

type (
	EventAdapter struct {
		ze *zerolog.Event
	}
)

func (e *EventAdapter) Caller(skip ...int) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return &EventAdapter{ze: e.ze.Caller(skip...)}
}

func (e *EventAdapter) Err(err error) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return &EventAdapter{ze: e.ze.Err(err)}
}

func (e *EventAdapter) Str(key, value string) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return &EventAdapter{ze: e.ze.Str(key, value)}
}

func (e *EventAdapter) Bytes(key string, value []byte) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return &EventAdapter{ze: e.ze.Bytes(key, value)}
}

func (e *EventAdapter) Int(key string, value int) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return &EventAdapter{ze: e.ze.Int(key, value)}
}

func (e *EventAdapter) Any(key string, value any) mrlog.LoggerEvent {
	if e == nil {
		return e
	}

	return &EventAdapter{ze: e.ze.Any(key, value)}
}

func (e *EventAdapter) Msg(message string) {
	if e == nil {
		return
	}

	e.ze.Msg(message)
}

func (e *EventAdapter) Msgf(format string, args ...any) {
	if e == nil {
		return
	}

	e.ze.Msgf(format, args...)
}

func (e *EventAdapter) MsgFunc(createMsg func() string) {
	if e == nil {
		return
	}

	e.ze.MsgFunc(createMsg)
}

func (e *EventAdapter) Send() {
	if e == nil {
		return
	}

	e.ze.Send()
}
