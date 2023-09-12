package mrenv

import (
    "context"

    "github.com/mondegor/go-webcore/mrcore"
)

type (
	ctxLogger struct{}
)

func WithLogger(ctx context.Context, value mrcore.Logger) context.Context {
    return context.WithValue(ctx, ctxLogger{}, value)
}

func Logger(ctx context.Context) mrcore.Logger {
    value, ok := ctx.Value(ctxLogger{}).(mrcore.Logger)

    if ok {
        return value
    }

    return mrcore.DefaultLogger()
}
