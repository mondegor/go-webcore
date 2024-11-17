package nopemitter

import (
	"context"
)

type (
	// Emitter - заглушка отправителя событий.
	Emitter struct{}
)

// New - создаёт объект Emitter.
func New() *Emitter {
	return &Emitter{}
}

// Emit - эмулирует отправку указанного события.
func (e *Emitter) Emit(_ context.Context, _ string, _ any) {
}
