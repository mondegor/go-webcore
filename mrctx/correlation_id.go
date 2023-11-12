package mrctx

import (
    "context"

    "github.com/google/uuid"
)

// https://www.rapid7.com/blog/post/2016/12/23/the-value-of-correlation-ids/

type (
    ctxCorrelationID struct{}
)

func WithCorrelationID(ctx context.Context, value string) context.Context {
    return context.WithValue(ctx, ctxCorrelationID{}, value)
}

func CorrelationID(ctx context.Context) string {
    value, ok := ctx.Value(ctxCorrelationID{}).(string)

    if ok {
        return value
    }

    return GenCorrelationID()
}

func GenCorrelationID() string {
    return uuid.New().String()
}
