package nopeventemitter

import (
	"context"

	"github.com/mondegor/go-webcore/mrsender"
)

type (
	// EventEmitter - comment struct.
	EventEmitter struct{}
)

// Make sure the Image conforms with the mrsender.EventEmitter interface.
var _ mrsender.EventEmitter = (*EventEmitter)(nil)

// New - создаёт объект EventEmitter.
func New() *EventEmitter {
	return &EventEmitter{}
}

// Emit - comment method.
func (e *EventEmitter) Emit(_ context.Context, _ string, _ any) {
}

// EmitWithSource - comment method.
func (e *EventEmitter) EmitWithSource(_ context.Context, _, _ string, _ any) {
}
