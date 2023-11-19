package mrctx

import "context"

// https://www.rapid7.com/blog/post/2016/12/23/the-value-of-correlation-ids/

func CorrelationID(ctx context.Context) string {
	value, ok := ctx.Value(ctxClientTools{}).(ClientTools)

	if ok {
		return value.CorrelationID
	}

	return "correlation-id-not-found"
}
