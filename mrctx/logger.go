package mrctx

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
)

func Logger(ctx context.Context) mrcore.Logger {
	value, ok := ctx.Value(ctxClientTools{}).(ClientTools)

	if ok {
		return value.Logger
	}

	return mrcore.DefaultLogger()
}
