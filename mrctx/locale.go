package mrenv

import (
    "context"

    "github.com/mondegor/go-sysmess/mrlang"
)

type (
	ctxLocale struct{}
)

func WithLocale(ctx context.Context, value *mrlang.Locale) context.Context {
    return context.WithValue(ctx, ctxLocale{}, value)
}

func Locale(ctx context.Context) *mrlang.Locale {
    value, ok := ctx.Value(ctxLocale{}).(*mrlang.Locale)

    if ok {
        return value
    }

    return mrlang.DefaultLocale()
}
