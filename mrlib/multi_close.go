package mrlib

import (
	"context"
	"fmt"
	"io"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

// CallEachFunc - comment func.
func CallEachFunc(ctx context.Context, list []func(ctx context.Context)) {
	for _, f := range list {
		f(ctx)
	}
}

// CloseFunc - comment func.
func CloseFunc(object io.Closer) func(ctx context.Context) {
	return func(ctx context.Context) {
		Close(ctx, object)
	}
}

// Close - comment func.
func Close(ctx context.Context, object io.Closer) {
	logger := mrlog.Ctx(ctx)

	if err := object.Close(); err != nil {
		logger.Error().
			Str("io.Closer", fmt.Sprintf("%#v", object)).
			Err(mrcore.ErrInternalFailedToClose.Wrap(err)).
			Send()
	} else {
		logger.Info().Msgf("Object %T was closed", object)
	}
}
