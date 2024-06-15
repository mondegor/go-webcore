package mrlogadapter

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
)

type (
	// EventEmitter - comment struct.
	EventEmitter struct {
		logger mrlog.Logger
	}
)

// Make sure the Image conforms with the mrsender.EventEmitter interface.
var _ mrsender.EventEmitter = (*EventEmitter)(nil)

// NewEventEmitter - создаёт объект EventEmitter.
func NewEventEmitter(logger mrlog.Logger) *EventEmitter {
	return &EventEmitter{
		logger: logger,
	}
}

// Emit - comment method.
func (l *EventEmitter) Emit(_ context.Context, eventName string, object any) {
	l.logger.Info().Str("event", eventName).Any("object", object).Send()
}

// EmitWithSource - comment method.
func (l *EventEmitter) EmitWithSource(_ context.Context, eventName, source string, object any) {
	l.logger.Info().Str("event", eventName).Str("source", source).Any("object", object).Send()
}
