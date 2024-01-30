package mrfactory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"
)

func NewAppSection(ctx context.Context, opts mrperms.AppSectionOptions, access mrperms.AccessControl) mrperms.AppSection {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msgf("Init section %s with root path '%s' and privilege '%s'", opts.Caption, opts.RootPath, opts.Privilege)
	logger.Debug().Msgf("secret=%s, audience: %s", opts.AuthSecret, opts.AuthAudience)

	if !access.HasPrivilege(opts.Privilege) {
		logger.Warn().Msgf(
			"Privilege '%s' is not registered for section '%s', perhaps, it is not registered in the config or is not associated with any role",
			opts.Privilege,
			opts.Caption,
		)
	}

	return mrperms.NewAppSection(opts)
}
