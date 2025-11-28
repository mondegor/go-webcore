package mraccess

type (
	// entryRole - роль с привязанными к ней привилегиями и разрешениями.
	entryRole struct {
		roleID uint16
		rights RightsSource
	}
)

// NewRole - создаёт объект  RightsChecker.
func NewRole(role string, rights RightsSource) RightsChecker {
	if roleIDs := rights.RoleIDsByNames([]string{role}); len(roleIDs) > 0 {
		return &entryRole{
			roleID: roleIDs[0],
			rights: rights,
		}
	}

	return &entryRole{
		rights: rights,
	}
}

// HasPrivilege - сообщает о наличии указанной привилегии.
func (r *entryRole) HasPrivilege(name string) bool {
	return r.roleInArray(
		r.rights.RoleIDsByPrivilege(name),
	)
}

// HasPermission - сообщает о наличии указанного разрешения.
func (r *entryRole) HasPermission(name string) bool {
	return r.roleInArray(
		r.rights.RoleIDsByPermission(name),
	)
}

func (r *entryRole) roleInArray(roleIDs []uint16) bool {
	for i := range roleIDs {
		if roleIDs[i] == r.roleID {
			return true
		}
	}

	return false
}
