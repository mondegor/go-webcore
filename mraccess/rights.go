package mraccess

type (
	// RightsAvailability - comment interface.
	RightsAvailability interface {
		HasPrivilege(name string) bool
		HasPermission(name string) bool
	}

	// RightsSource - comment interface.
	RightsSource interface {
		RoleIDsByNames(roles []string) []uint16
		RoleIDsByPrivilege(name string) []uint16
		RoleIDsByPermission(name string) []uint16
		RightsAvailability
	}

	// RightsChecker - проверяет права доступа к указанным привилегиям и разрешениям.
	RightsChecker interface {
		CheckPrivilege(name string) bool
		CheckPermission(name string) bool
	}

	// RightsGetter - comment interface.
	RightsGetter interface {
		Rights(name string) RightsChecker
	}
)
