package mrsender

import (
	"context"
)

type (
	nopEventEmitter struct{}
)

func NewNopEventEmitter() EventEmitter {
	return &nopEventEmitter{}
}

func (e *nopEventEmitter) Emit(ctx context.Context, eventName string, object any) {
}

func (e *nopEventEmitter) EmitWithSource(ctx context.Context, eventName, source string, object any) {
}
