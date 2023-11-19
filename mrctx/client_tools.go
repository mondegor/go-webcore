package mrctx

import (
	"context"
	"fmt"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
)

type (
	ctxClientTools struct{}

	ClientTools struct {
		CorrelationID string
		Logger        mrcore.Logger
		Locale        *mrlang.Locale
		Validator     mrcore.Validator
	}
)

func WithClientTools(
	ctx context.Context,
	correlationID string,
	logger mrcore.Logger,
	locale *mrlang.Locale,
	validator mrcore.Validator,
) context.Context {
	return context.WithValue(
		ctx,
		ctxClientTools{},
		ClientTools{
			CorrelationID: correlationID,
			Logger:        logger,
			Locale:        locale,
			Validator:     validator,
		},
	)
}

func GetClientTools(ctx context.Context) (ClientTools, error) {
	value, ok := ctx.Value(ctxClientTools{}).(ClientTools)

	if ok {
		return value, nil
	}

	return ClientTools{}, fmt.Errorf("client request tools not found in context")
}
