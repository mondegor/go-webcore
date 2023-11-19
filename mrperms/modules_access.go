package mrperms

import (
	"fmt"
)

type (
	roleMap       map[string]int32
	privilegeMap  map[string][]int32
	permissionMap map[string][]int32

	ModulesAccess struct {
		roles        roleMap       // map to rolesIDs
		privileges   privilegeMap  // map to roles
		permissions  permissionMap // map to roles
		guestsRole   string
		guestsRoleID int32
	}

	ModulesAccessOptions struct {
		RolesDirPath  string
		RolesFileType string
		Roles         []string
		Privileges    []string
		Permissions   []string
		GuestRole     string // optional
	}
)

func NewModulesAccess(opt ModulesAccessOptions) (*ModulesAccess, error) {
	if len(opt.Roles) == 0 {
		return nil, fmt.Errorf("opt.Roles is required")
	}

	if opt.GuestRole == "" {
		opt.GuestRole = opt.Roles[0]
	} else if !roleInArray(opt.GuestRole, opt.Roles) {
		return nil, fmt.Errorf("opt.GuestRole='%s' not found in opt.Roles", opt.GuestRole)
	}

	ma := ModulesAccess{
		roles:       make(roleMap, len(opt.Roles)),
		privileges:  make(privilegeMap, 0),
		permissions: make(permissionMap, 0),
	}

	var roleID int32

	for pos, roleName := range opt.Roles {
		if _, ok := ma.roles[roleName]; ok {
			return nil, fmt.Errorf("duplicate role %s detected (pos: %d)", roleName, pos+1)
		}

		roleID++
		ma.roles[roleName] = roleID

		config, err := loadRoleConfig(roleName, opt.RolesDirPath, opt.RolesFileType)

		if err != nil {
			return nil, err
		}

		for _, priv := range config.Privileges {
			if !stringInArray(priv, opt.Privileges) {
				return nil, fmt.Errorf("privilege '%s' is not registered in opt.Privileges", priv)
			}

			ma.privileges[priv] = append(ma.privileges[priv], roleID)
		}

		for _, perm := range config.Permissions {
			if !stringInArray(perm, opt.Permissions) {
				return nil, fmt.Errorf("permission '%s' is not registered in opt.Permissions", perm)
			}

			ma.permissions[perm] = append(ma.permissions[perm], roleID)
		}

		if ma.guestsRoleID == 0 && opt.GuestRole == roleName {
			ma.guestsRole = roleName
			ma.guestsRoleID = roleID
		}
	}

	return &ma, nil
}

func (a *ModulesAccess) NewRoleGroup(roles []string) *RoleGroup {
	return newRoleGroup(a, roles)
}

func (a *ModulesAccess) GuestRole() string {
	return a.guestsRole
}

func (a *ModulesAccess) CheckPrivilege(rolesIDs []int32, name string) bool {
	privRolesIDs, ok := a.privileges[name]

	if !ok {
		return false
	}

	return isArraysIntersection(rolesIDs, privRolesIDs)
}

func (a *ModulesAccess) CheckPermission(rolesIDs []int32, name string) bool {
	permRolesIDs, ok := a.permissions[name]

	if !ok {
		return false
	}

	return isArraysIntersection(rolesIDs, permRolesIDs)
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

func (a *ModulesAccess) roleNamesToIDs(roles []string) []int32 {
	var roleIDs []int32

	for _, role := range roles {
		if id, ok := a.roles[role]; ok {
			roleIDs = append(roleIDs, id)
		}
	}

	return roleIDs
}

func isArraysIntersection(ids1, ids2 []int32) bool {
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

func roleInArray(roleName string, roleNames []string) bool {
	for i := range roleNames {
		if roleName == roleNames[i] {
			return true
		}
	}

	return false
}
