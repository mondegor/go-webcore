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
	// Используется как центральный источник информации о правах доступа.
	RightsSource interface {
		RightsChecker

		// RoleIDsByNames - возвращает ID ролей по их именам.
		RoleIDsByNames(roles []string) []uint16

		// RoleIDsByPrivilege - возвращает ID всех ролей, обладающих указанной привилегией.
		RoleIDsByPrivilege(name string) []uint16

		// RoleIDsByPermission - возвращает ID всех ролей, обладающих указанным разрешением.
		RoleIDsByPermission(name string) []uint16
	}

	// RightsGetter - возвращает объект для проверки привилегий
	// и разрешений для указанной группы.
	RightsGetter interface {
		// Rights - возвращает RightsChecker для указанной группы.
		Rights(group string) RightsChecker
	}
)
