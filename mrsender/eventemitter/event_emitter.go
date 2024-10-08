package eventemitter

import (
	"context"

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
func (l *Emitter) Emit(ctx context.Context, eventName string, object any) {
	for _, r := range l.receivers {
		r.Receive(ctx, eventName, "none", object)
	}
}

// EmitWithSource - отправляет указанное событие включающее источник.
func (l *Emitter) EmitWithSource(ctx context.Context, eventName, source string, object any) {
	for _, r := range l.receivers {
		r.Receive(ctx, eventName, source, object)
	}
}
