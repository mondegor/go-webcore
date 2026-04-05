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
	// RightsChecker - проверяет наличие указанных привилегий и разрешений
	// у сущности (пользователя, роли и т.д.).
	RightsChecker interface {
		HasPrivilege(name string) bool
		HasPermission(name string) bool
	}

	// RightsSource - предоставляет методы для получения ролей,
	// разрешений и привилегий, а также их проверки.
	RightsSource interface {
		RightsChecker

		RoleIDsByNames(roles []string) []uint16
		RoleIDsByPrivilege(name string) []uint16
		RoleIDsByPermission(name string) []uint16
	}

	// RightsGetter - возвращает объект для проверки привилегий
	// и разрешений для указанной группы.
	RightsGetter interface {
		Rights(group string) RightsChecker
	}
)
