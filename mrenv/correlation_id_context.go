package mrenv

import (
    "context"

    "github.com/google/uuid"
)

type (
	ctxCorrelationId struct{}
)

func WithCorrelationId(ctx context.Context, value string) context.Context {
    return context.WithValue(ctx, ctxCorrelationId{}, value)
}

func CorrelationIdFromContext(ctx context.Context) string {
    value, ok := ctx.Value(ctxCorrelationId{}).(string)

    if ok {
        return value
    }

    return GenCorrelationId()
}

func GenCorrelationId() string {
    return uuid.New().String()
}
