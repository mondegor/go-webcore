package wire

import (
	"strings"

	"github.com/mondegor/go-sysmess/mraccess/provider/filestorage"
	"github.com/mondegor/go-sysmess/mrlog"
)

// InitPermsProvider - создаёт и инициализирует провайдер ролей и разрешений из файлового хранилища.
// Параметры
//   - rolesDirPath - путь к директории с файлами ролей;
//   - roles - список ролей;
//   - allowedPrivileges, allowedPermissions - списки разрешенных привилегий и разрешений для регистрации.
func InitPermsProvider(
	logger mrlog.Logger,
	rolesDirPath string,
	roles []string,
	allowedPrivileges []string,
	allowedPermissions []string,
) (*filestorage.PermsProvider, error) {
	mrlog.Info(logger, "Create and init roles and permissions for app")

	provider, err := filestorage.NewPermsProvider(
		rolesDirPath,
		roles,
		filestorage.WithAllowedPrivileges(allowedPrivileges),
		filestorage.WithAllowedPermissions(allowedPermissions),
	)
	if err != nil {
		return nil, err
	}

	info := filestorage.ExtractProviderInfo(provider)

	mrlog.Info(logger, "Registered roles: "+strings.Join(info.Roles, ", "))
	mrlog.Info(logger, "Registered privileges: "+strings.Join(info.Privileges, ", "))
	mrlog.Debug(logger, "Registered permissions:\n - "+strings.Join(info.Permissions, ",\n - "))

	return provider, nil
}
