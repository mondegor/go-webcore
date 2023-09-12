package mrctx

import (
    "context"
)

func WithEnumList(ctx context.Context, ctxKey any, items []string) context.Context {
    return context.WithValue(ctx, ctxKey, items)
}

func EnumList(ctx context.Context, ctxKey any) []string {
    value, ok := ctx.Value(ctxKey).([]string)

    if ok {
        return value
    }

    return []string{}
}
