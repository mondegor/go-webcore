package initing

import (
	"strings"

	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mraccess/config/filestorage"
)

// InitPermsProvider - создаёт объект filestorage.PermsProvider.
func InitPermsProvider(
	logger mrlog.Logger,
	rolesDirPath string,
	roles []string,
	privileges []string,
	permissions []string,
) (*filestorage.PermsProvider, error) {
	mrlog.Info(logger, "Create and init roles and permissions for app")

	provider, err := filestorage.NewPermsProvider(
		rolesDirPath,
		roles,
		filestorage.WithPrivileges(privileges),
		filestorage.WithPermissions(permissions),
	)
	if err != nil {
		return nil, err
	}

	info := filestorage.NewRegisteredPermsInfo(provider)

	mrlog.Info(logger, "Registered roles: "+strings.Join(info.Roles, ", "))
	mrlog.Info(logger, "Registered privileges: "+strings.Join(info.Privileges, ", "))
	mrlog.Debug(logger, "Registered permissions:\n - "+strings.Join(info.Permissions, ",\n - "))

	return provider, nil
}
