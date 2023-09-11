package mrenv

import (
    "context"
    "net"
)

type (
	ctxUserIP struct{}
)

func WithUserIp(ctx context.Context, value net.IP) context.Context {
    return context.WithValue(ctx, ctxUserIP{}, value)
}

func UserIpFromContext(ctx context.Context) net.IP {
    value, ok := ctx.Value(ctxUserIP{}).(net.IP)

    if ok {
        return value
    }

    return net.IP{}
}
