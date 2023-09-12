package mrenv

import "context"

func WithInt64List(ctx context.Context, ctxKey any, items []int64) context.Context {
    return context.WithValue(ctx, ctxKey, items)
}

func Int64List(ctx context.Context, ctxKey any) []int64 {
    value, ok := ctx.Value(ctxKey).([]int64)

    if ok {
        return value
    }

    return []int64{}
}
