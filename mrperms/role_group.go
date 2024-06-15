package mrperms

type (
	roleGroup struct {
		roleIDs []uint16
		access  *RoleAccessControl
	}
)

func newRoleGroup(access *RoleAccessControl, roles []string) *roleGroup {
	return &roleGroup{
		roleIDs: access.roleNamesToIDs(roles),
		access:  access,
	}
}

// IsGuestAccess - если ролей нет или присутствует ровно одна роль и она гостевая.
func (g *roleGroup) IsGuestAccess() bool {
	return len(g.roleIDs) == 0 || len(g.roleIDs) == 1 && g.roleIDs[0] == g.access.guestsRoleID
}

// CheckPrivilege - comment method.
func (g *roleGroup) CheckPrivilege(name string) bool {
	return g.access.CheckPrivilege(g.roleIDs, name)
}

// CheckPermission - comment method.
func (g *roleGroup) CheckPermission(name string) bool {
	return g.access.CheckPermission(g.roleIDs, name)
}
