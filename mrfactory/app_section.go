package mrfactory

import (
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrperms"
)

func NewAppSection(opt mrperms.AppSectionOptions, access mrcore.AccessControl, logger mrcore.Logger) *mrperms.AppSection {
	logger.Info("Init section %s with root path '%s' and privilege '%s'", opt.Caption, opt.RootPath, opt.Privilege)
	logger.Debug("secret=%s, audience: %s", opt.AuthSecret, opt.AuthAudience)

	if !access.HasPrivilege(opt.Privilege) {
		logger.Warning(
			"privilege '%s' is not registered for section '%s', perhaps, it is not registered in the config or is not associated with any role",
			opt.Privilege,
			opt.Caption,
		)
	}

	return mrperms.NewAppSection(opt)
}
