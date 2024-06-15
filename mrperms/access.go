package mrperms

type (
	// AccessRightsFactory - выдаёт объект AccessRights, который проверяет права доступа для указанных members.
	AccessRightsFactory interface {
		NewAccessRights(members ...string) AccessRights
	}

	// AccessRights - проверка прав доступа конкретного объекта к привилегиям и разрешениям.
	AccessRights interface {
		IsGuestAccess() bool
		CheckPrivilege(name string) bool
		CheckPermission(name string) bool
	}
)
