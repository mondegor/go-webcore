package mrctx

import (
	"context"

	"github.com/mondegor/go-sysmess/mrlang"
)

func Locale(ctx context.Context) *mrlang.Locale {
	value, ok := ctx.Value(ctxClientTools{}).(ClientTools)

	if ok {
		return value.Locale
	}

	return mrlang.DefaultLocale()
}
