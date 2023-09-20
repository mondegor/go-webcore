package mrctx

import (
    "context"

    "github.com/google/uuid"
)

// https://www.rapid7.com/blog/post/2016/12/23/the-value-of-correlation-ids/

type (
    ctxCorrelationId struct{}
)

func WithCorrelationId(ctx context.Context, value string) context.Context {
    return context.WithValue(ctx, ctxCorrelationId{}, value)
}

func CorrelationId(ctx context.Context) string {
    value, ok := ctx.Value(ctxCorrelationId{}).(string)

    if ok {
        return value
    }

    return GenCorrelationId()
}

func GenCorrelationId() string {
    return uuid.New().String()
}
