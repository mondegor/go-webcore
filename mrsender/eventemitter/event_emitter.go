package eventemitter

import (
	"context"
	"strings"

	"github.com/mondegor/go-webcore/mrsender"
)

type (
	// Emitter - синхронный отправитель событий.
	Emitter struct {
		receivers []mrsender.EventReceiver
	}
)

// New - создаёт объект Emitter.
func New(receivers ...mrsender.EventReceiver) *Emitter {
	return &Emitter{
		receivers: receivers,
	}
}

// Emit - отправляет указанное событие.
func (e *Emitter) Emit(ctx context.Context, eventName string, object any) {
	source := mrsender.DefaultSource

	if s, n, ok := strings.Cut(eventName, mrsender.SourceEventSeparator); ok {
		source = s
		eventName = n
	}

	for _, r := range e.receivers {
		r.Receive(ctx, eventName, source, object)
	}
}
