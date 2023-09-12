package mrenv

import "context"

func WithInt64(ctx context.Context, ctxKey any, value int64) context.Context {
    return context.WithValue(ctx, ctxKey, value)
}

func Int64(ctx context.Context, ctxKey any) int64 {
    value, ok := ctx.Value(ctxKey).(int64)

    if ok {
        return value
    }

    return 0
}
