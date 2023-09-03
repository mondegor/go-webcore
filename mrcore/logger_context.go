package mrcore

import (
    "context"

    "github.com/mondegor/go-core/mrlog"
)

type ctxLogger struct{}

func WithLogger(ctx context.Context, logger mrlog.Logger) context.Context {
    return context.WithValue(ctx, ctxLogger{}, logger)
}

func ExtractLogger(ctx context.Context) mrlog.Logger {
    value, ok := ctx.Value(ctxLogger{}).(mrlog.Logger)

    if ok {
        return value
    }

    return mrlog.Default()
}
