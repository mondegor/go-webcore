package mraccess

type (
	// RoleGroup - группа с привязанными к ней ролями.
	RoleGroup struct {
		Name  string
		Roles []string
	}

	// entryRoleGroup - группа ролей с привязанными к ним привилегиями и разрешениями.
	entryRoleGroup struct {
		roleIDs []uint16
		rights  RightsSource
	}
)

// NewRoleGroup - создаёт объект RightsChecker.
func NewRoleGroup(roles []string, rights RightsSource) RightsChecker {
	return &entryRoleGroup{
		roleIDs: rights.RoleIDsByNames(roles),
		rights:  rights,
	}
}

// HasPrivilege - сообщает о наличии указанной привилегии.
func (g *entryRoleGroup) HasPrivilege(name string) bool {
	return g.isIntersection(
		g.rights.RoleIDsByPrivilege(name),
	)
}

// HasPermission - сообщает о наличии указанного разрешения.
func (g *entryRoleGroup) HasPermission(name string) bool {
	return g.isIntersection(
		g.rights.RoleIDsByPermission(name),
	)
}

func (g *entryRoleGroup) isIntersection(roleIDs []uint16) bool {
	for i := range g.roleIDs {
		for j := range roleIDs {
			if g.roleIDs[i] == roleIDs[j] {
				return true
			}
		}
	}

	return false
}
