package mrenv

import (
    "context"
)

const (
    PlatformMobile = "MOBILE"
    PlatformWeb = "WEB"
)

type (
	ctxPlatform struct{}
)

func WithPlatform(ctx context.Context, value string) context.Context {
    return context.WithValue(ctx, ctxPlatform{}, value)
}

func PlatformFromContext(ctx context.Context) string {
    value, ok := ctx.Value(ctxPlatform{}).(string)

    if ok {
        return value
    }

    return PlatformWeb
}
