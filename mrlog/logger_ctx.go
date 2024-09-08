package mrlog

import (
	"context"
)

type (
	ctxKey struct{}
)

// WithContext - возвращается контекст с указанным логгером.
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// Ctx - возвращается логер из указанного контекста,
// если он не был установлен ранее, то возвращается логгер по умолчанию.
func Ctx(ctx context.Context) Logger {
	if value, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return value
	}

	return def
}

// HasCtx - проверяется, что логгер явно содержится в контексте.
func HasCtx(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKey{}).(Logger)

	return ok
}
