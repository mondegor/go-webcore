package mraccess

const (
	// PrivilegePublic - привилегия для всех.
	PrivilegePublic = "public"

	// PermissionAnyUser - разрешение для любого пользователя.
	PermissionAnyUser = "any-user"

	// PermissionGuestOnly - разрешение только для гостя.
	PermissionGuestOnly = "guest-only"
)

type (
	// RightsChecker - проверяет наличие указанных привилегий и разрешений сущности.
	RightsChecker interface {
		HasPrivilege(name string) bool
		HasPermission(name string) bool
	}

	// RightsSource - интерфейс источника ролей, разрешений и привилегий.
	RightsSource interface {
		RoleIDsByNames(roles []string) []uint16
		RoleIDsByPrivilege(name string) []uint16
		RoleIDsByPermission(name string) []uint16
		RightsChecker
	}

	// RightsGetter - возвращает объект проверяющий наличие указанных привилегий и разрешений сущности.
	RightsGetter interface {
		Rights(group string) RightsChecker
	}
)
