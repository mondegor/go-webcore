package mrperms

type (
	RoleGroup struct {
		roleIDs []uint16
		access  *AccessControl
	}
)

func newRoleGroup(access *AccessControl, roles []string) *RoleGroup {
	return &RoleGroup{
		roleIDs: access.roleNamesToIDs(roles),
		access:  access,
	}
}

// IsGuestAccess - если ролей нет или присутствует ровно одна роль и она гостевая
func (g *RoleGroup) IsGuestAccess() bool {
	return len(g.roleIDs) == 0 || len(g.roleIDs) == 1 && g.roleIDs[0] == g.access.guestsRoleID
}

func (g *RoleGroup) CheckPrivilege(name string) bool {
	return g.access.CheckPrivilege(g.roleIDs, name)
}

func (g *RoleGroup) CheckPermission(name string) bool {
	return g.access.CheckPermission(g.roleIDs, name)
}
