package mrfactory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
)

// InfoCreateModule - comment func.
func InfoCreateModule(ctx context.Context, name string) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init module '%s'", name)
}

// InfoCreateUnit - comment func.
func InfoCreateUnit(ctx context.Context, name string) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init unit '%s' of the module", name)
}
