package mrenv

import (
    "context"

    "github.com/mondegor/go-webcore/mrcore"
)

type (
	ctxPlatform struct{}
)

func WithPlatform(ctx context.Context, value string) context.Context {
    return context.WithValue(ctx, ctxPlatform{}, value)
}

func Platform(ctx context.Context) string {
    value, ok := ctx.Value(ctxPlatform{}).(string)

    if ok {
        return value
    }

    return mrcore.PlatformWeb
}
