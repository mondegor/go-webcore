package mrperms

import (
	"fmt"
)

type (
	roleMap       map[string]uint16
	privilegeMap  map[string][]uint16
	permissionMap map[string][]uint16

	RoleAccessControl struct {
		roles        roleMap       // map to rolesIDs
		privileges   privilegeMap  // map to roles
		permissions  permissionMap // map to roles
		guestsRole   string
		guestsRoleID uint16
	}

	AccessControlOptions struct {
		RolesDirPath  string
		RolesFileType string
		Roles         []string
		Privileges    []string
		Permissions   []string
		GuestRole     string // optional
	}
)

func NewAccessControl(opts AccessControlOptions) (*RoleAccessControl, error) {
	if len(opts.Roles) == 0 {
		return nil, fmt.Errorf("opts.Roles is required")
	}

	if opts.GuestRole == "" {
		opts.GuestRole = opts.Roles[0]
	} else if !roleInArray(opts.GuestRole, opts.Roles) {
		return nil, fmt.Errorf("opts.GuestRole='%s' not found in opts.Roles", opts.GuestRole)
	}

	ma := RoleAccessControl{
		roles:       make(roleMap, len(opts.Roles)),
		privileges:  make(privilegeMap, 0),
		permissions: make(permissionMap, 0),
	}

	var roleID uint16

	for pos, roleName := range opts.Roles {
		if _, ok := ma.roles[roleName]; ok {
			return nil, fmt.Errorf("duplicate role %s detected (pos: %d)", roleName, pos+1)
		}

		roleID++
		ma.roles[roleName] = roleID

		config, err := loadRoleConfig(roleName, getFilePath(opts.RolesDirPath, roleName))

		if err != nil {
			return nil, err
		}

		for _, priv := range config.Privileges {
			if !stringInArray(priv, opts.Privileges) {
				return nil, fmt.Errorf("privilege '%s' is not registered in opts.Privileges", priv)
			}

			ma.privileges[priv] = append(ma.privileges[priv], roleID)
		}

		for _, perm := range config.Permissions {
			if !stringInArray(perm, opts.Permissions) {
				return nil, fmt.Errorf("permission '%s' is not registered in opts.Permissions", perm)
			}

			ma.permissions[perm] = append(ma.permissions[perm], roleID)
		}

		if ma.guestsRoleID == 0 && opts.GuestRole == roleName {
			ma.guestsRole = roleName
			ma.guestsRoleID = roleID
		}
	}

	return &ma, nil
}

func (a *RoleAccessControl) NewAccessRights(roles ...string) AccessRights {
	return newRoleGroup(a, roles)
}

func (a *RoleAccessControl) GuestRole() string {
	return a.guestsRole
}

func (a *RoleAccessControl) HasPrivilege(name string) bool {
	_, ok := a.privileges[name]

	return ok
}

func (a *RoleAccessControl) CheckPrivilege(rolesIDs []uint16, name string) bool {
	privRolesIDs, ok := a.privileges[name]

	if !ok {
		return false
	}

	return isArraysIntersection(rolesIDs, privRolesIDs)
}

func (a *RoleAccessControl) HasPermission(name string) bool {
	_, ok := a.permissions[name]

	return ok
}

func (a *RoleAccessControl) CheckPermission(rolesIDs []uint16, name string) bool {
	permRolesIDs, ok := a.permissions[name]

	if !ok {
		return false
	}

	return isArraysIntersection(rolesIDs, permRolesIDs)
}

func (a *RoleAccessControl) RegisteredRoles() []string {
	roles := make([]string, len(a.roles))
	i := 0

	for name := range a.roles {
		roles[i] = name
		i++
	}

	return roles
}

func (a *RoleAccessControl) RegisteredPrivileges() []string {
	privileges := make([]string, len(a.privileges))
	i := 0

	for name := range a.privileges {
		privileges[i] = name
		i++
	}

	return privileges
}

func (a *RoleAccessControl) RegisteredPermissions() []string {
	permissions := make([]string, len(a.permissions))
	i := 0

	for name := range a.permissions {
		permissions[i] = name
		i++
	}

	return permissions
}

func (a *RoleAccessControl) roleNamesToIDs(roles []string) []uint16 {
	var roleIDs []uint16

	for _, role := range roles {
		if id, ok := a.roles[role]; ok {
			roleIDs = append(roleIDs, id)
		}
	}

	return roleIDs
}

func isArraysIntersection(ids1, ids2 []uint16) bool {
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
