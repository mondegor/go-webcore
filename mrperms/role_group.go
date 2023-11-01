package mrperms

type (
    RoleGroup struct {
        roleIds []int32
        access *ModulesAccess
    }
)

func newRoleGroup(access *ModulesAccess, roles []string) *RoleGroup {
    return &RoleGroup{
        roleIds: access.roleNamesToIds(roles),
        access: access,
    }
}

func (g *RoleGroup) CheckPrivilege(name string) bool {
    return g.access.CheckPrivilege(g.roleIds, name)
}

func (g *RoleGroup) CheckPermission(name string) bool {
    return g.access.CheckPermission(g.roleIds, name)
}
