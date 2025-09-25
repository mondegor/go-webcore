package filestorage

import (
	"errors"
	"fmt"
)

type (
	// PermsProvider - comment struct.
	PermsProvider struct {
		roles       map[string]uint16   // map to rolesIDs
		privileges  map[string][]uint16 // map to roles
		permissions map[string][]uint16 // map to roles
	}

	permsProvider struct {
		privileges  []string
		permissions []string
	}
)

// NewPermsProvider - создаёт объект PermsProvider.
func NewPermsProvider(dirPath string, roles []string, opts ...Option) (*PermsProvider, error) {
	if len(roles) == 0 {
		return nil, errors.New("roles is required")
	}

	p := &PermsProvider{
		roles:       make(map[string]uint16, len(roles)),
		privileges:  make(map[string][]uint16),
		permissions: make(map[string][]uint16),
	}

	cfg := permsProvider{}

	for _, opt := range opts {
		opt(&cfg)
	}

	var roleID uint16

	for i, roleName := range roles {
		if _, ok := p.roles[roleName]; ok {
			return nil, fmt.Errorf("duplicate role %s detected in param 'roles' (pos: %d)", roleName, i+1)
		}

		roleID++
		p.roles[roleName] = roleID

		fileCfg, err := loadRoleConfig(getFilePath(dirPath, roleName))
		if err != nil {
			return nil, err
		}

		for _, priv := range fileCfg.Privileges {
			if !stringInArray(priv, cfg.privileges) {
				return nil, fmt.Errorf("privilege '%s' is not registered in opt.Privileges", priv)
			}

			p.privileges[priv] = append(p.privileges[priv], roleID)
		}

		for _, perm := range fileCfg.Permissions {
			if !stringInArray(perm, cfg.permissions) {
				return nil, fmt.Errorf("permission '%s' is not registered in opt.Permissions", perm)
			}

			p.permissions[perm] = append(p.permissions[perm], roleID)
		}
	}

	return p, nil
}

// RoleIDsByNames - comment method.
func (p *PermsProvider) RoleIDsByNames(roles []string) []uint16 {
	if len(roles) == 0 {
		return nil
	}

	roleIDs := make([]uint16, 0, len(roles))

	for i := range roles {
		if id, ok := p.roles[roles[i]]; ok {
			roleIDs = append(roleIDs, id)
		}
	}

	return roleIDs
}

// HasPrivilege - comment method.
func (p *PermsProvider) HasPrivilege(name string) bool {
	_, ok := p.privileges[name]

	return ok
}

// RoleIDsByPrivilege - comment method.
func (p *PermsProvider) RoleIDsByPrivilege(name string) []uint16 {
	return p.privileges[name]
}

// HasPermission - comment method.
func (p *PermsProvider) HasPermission(name string) bool {
	_, ok := p.permissions[name]

	return ok
}

// RoleIDsByPermission - comment method.
func (p *PermsProvider) RoleIDsByPermission(name string) []uint16 {
	return p.permissions[name]
}

func stringInArray(value string, values []string) bool {
	for i := range values {
		if value == values[i] {
			return true
		}
	}

	return false
}
