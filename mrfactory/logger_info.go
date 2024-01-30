package mrfactory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
)

func InfoCreateModule(ctx context.Context, name string) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init module '%s'", name)
}

func InfoCreateUnit(ctx context.Context, name string) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init unit '%s' of the module", name)
}
