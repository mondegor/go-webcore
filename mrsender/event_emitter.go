package mrsender

import (
	"context"
)

type (
	EventEmitter interface {
		Emit(ctx context.Context, eventName string, object any)
		EmitWithSource(ctx context.Context, eventName, source string, object any)
	}
)
