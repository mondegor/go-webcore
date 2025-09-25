package mrsender

import (
	"context"
)

type (
	// Tracer - comment interface.
	Tracer interface {
		Enabled() bool
		Trace(ctx context.Context, args ...any)
	}

	// TracerFunc - comment func type.
	TracerFunc func(ctx context.Context, args ...any)
)

// Trace - comment method.
func (f TracerFunc) Trace(ctx context.Context, args ...any) {
	f(ctx, args...)
}
