package mrperms

import (
    "fmt"
)

type (
    roleMap map[string]int32
    privilegeMap map[string][]int32
    permissionMap map[string][]int32

    ModulesAccess struct {
        roles roleMap // map to rolesIds
        privileges privilegeMap // map to roles
        permissions permissionMap // map to roles
        defaultRole string
        defaultRoleId int32
    }

    ModulesAccessOptions struct {
        RolesDirPath string
        RolesFileType string
        Roles []string
        Privileges []string
        Permissions []string
        DefaultRole string // optional
    }
)

func NewModulesAccess(opt ModulesAccessOptions) (*ModulesAccess, error) {
    if len(opt.Roles) == 0 {
        return nil, fmt.Errorf("opt.Roles is required")
    }

    if opt.DefaultRole == "" {
        opt.DefaultRole = opt.Roles[0]
    } else if !defaultRoleInArray(opt.DefaultRole, opt.Roles) {
        return nil, fmt.Errorf("opt.DefaultRole='%s' not found in opt.Roles", opt.DefaultRole)
    }

    ma := ModulesAccess{
        roles: make(roleMap, len(opt.Roles)),
        privileges: make(privilegeMap, 0),
        permissions: make(permissionMap, 0),
    }

    var roleId int32

    for pos, roleName := range opt.Roles {
        if _, ok := ma.roles[roleName]; ok {
            return nil, fmt.Errorf("duplicate role %s detected (pos: %d)", roleName, pos + 1)
        }

        roleId++
        ma.roles[roleName] = roleId

        config, err := loadRoleConfig(roleName, opt.RolesDirPath, opt.RolesFileType)

        if err != nil {
            return nil, err
        }

        for _, priv := range config.Privileges {
            if !stringInArray(priv, opt.Privileges) {
                return nil, fmt.Errorf("privilege '%s' is not registered in opt.Privileges", priv)
            }

            ma.privileges[priv] = append(ma.privileges[priv], roleId)
        }

        for _, perm := range config.Permissions {
            if !stringInArray(perm, opt.Permissions) {
                return nil, fmt.Errorf("permission '%s' is not registered in opt.Permissions", perm)
            }

            ma.permissions[perm] = append(ma.permissions[perm], roleId)
        }

        if opt.DefaultRole == roleName {
            ma.defaultRole = roleName
            ma.defaultRoleId = roleId
        }
    }

    return &ma, nil
}

func (a *ModulesAccess) NewRoleGroup(roles []string) *RoleGroup {
    return newRoleGroup(a, roles)
}

func (a *ModulesAccess) DefaultRole() string {
    return a.defaultRole
}

func (a *ModulesAccess) CheckPrivilege(rolesIds []int32, name string) bool {
    privRolesIds, ok := a.privileges[name]

    if !ok {
        return false
    }

    return isArraysIntersection(rolesIds, privRolesIds)
}

func (a *ModulesAccess) CheckPermission(rolesIds []int32, name string) bool {
    permRolesIds, ok := a.permissions[name]

    if !ok {
        return false
    }

    return isArraysIntersection(rolesIds, permRolesIds)
}

func (a *ModulesAccess) RegisteredRoles() []string {
    roles := make([]string, len(a.roles))
    i := 0

    for name := range a.roles {
        roles[i] = name
        i++
    }

    return roles
}

func (a *ModulesAccess) RegisteredPrivileges() []string {
    privileges := make([]string, len(a.privileges))
    i := 0

    for name := range a.privileges {
        privileges[i] = name
        i++
    }

    return privileges
}

func (a *ModulesAccess) RegisteredPermissions() []string {
    permissions := make([]string, len(a.permissions))
    i := 0

    for name := range a.permissions {
        permissions[i] = name
        i++
    }

    return permissions
}

func (a *ModulesAccess) roleNamesToIds(roles []string) []int32 {
    var roleIds []int32

    for _, role := range roles {
        if id, ok := a.roles[role]; ok {
            roleIds = append(roleIds, id)
        }
    }

    return roleIds
}

func isArraysIntersection(ids1 []int32, ids2 []int32) bool {
    for id1 := range ids1 {
        for id2 := range ids2 {
            if id1 == id2 {
                return true
            }
        }
    }

    return false
}

func stringInArray(value string, values []string) bool {
    for i := range values {
        if value == values[i] {
            return true
        }
    }

    return false
}

func defaultRoleInArray(roleName string, roleNames []string) bool {
    for i := range roleNames {
        if roleName == roleNames[i] {
            return true
        }
    }

    return false
}
