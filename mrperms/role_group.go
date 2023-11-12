package mrperms

type (
	RoleGroup struct {
		roleIDs []int32
		access  *ModulesAccess
	}
)

func newRoleGroup(access *ModulesAccess, roles []string) *RoleGroup {
	return &RoleGroup{
		roleIDs: access.roleNamesToIDs(roles),
		access:  access,
	}
}

func (g *RoleGroup) CheckPrivilege(name string) bool {
	return g.access.CheckPrivilege(g.roleIDs, name)
}

func (g *RoleGroup) CheckPermission(name string) bool {
	return g.access.CheckPermission(g.roleIDs, name)
}
