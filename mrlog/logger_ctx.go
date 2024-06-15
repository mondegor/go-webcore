package mrlog

import (
	"context"
)

type (
	ctxKey struct{}
)

// WithContext  - comment func.
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// Ctx  - comment func.
func Ctx(ctx context.Context) Logger {
	if value, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return value
	}

	return Default()
}
