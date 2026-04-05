package filestorage

import (
	"errors"
	"fmt"
)

type (
	// PermsProvider - управляет ролями, привилегиями и разрешениями из файлового хранилища.
	PermsProvider struct {
		roles       map[string]uint16   // map role names to role IDs
		privileges  map[string][]uint16 // map privilege names to role IDs
		permissions map[string][]uint16 // map permission names to role IDs
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
			return nil, fmt.Errorf("duplicate role detected in param 'roles' (role='%s', pos=%d)", roleName, i+1)
		}

		roleID++
		p.roles[roleName] = roleID

		fileCfg, err := loadRoleConfig(getFilePath(dirPath, roleName))
		if err != nil {
			return nil, err
		}

		for _, priv := range fileCfg.Privileges {
			if !stringInArray(priv, cfg.privileges) {
				return nil, fmt.Errorf("privilege is not registered in opt.Privileges (name='%s')", priv)
			}

			p.privileges[priv] = append(p.privileges[priv], roleID)
		}

		for _, perm := range fileCfg.Permissions {
			if !stringInArray(perm, cfg.permissions) {
				return nil, fmt.Errorf("permission is not registered in opt.Permissions (name='%s')", perm)
			}

			p.permissions[perm] = append(p.permissions[perm], roleID)
		}
	}

	return p, nil
}

// RoleIDsByNames - возвращает ID ролей по их названиям.
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

// HasPrivilege - сообщает о наличии указанной привилегии.
func (p *PermsProvider) HasPrivilege(name string) bool {
	_, ok := p.privileges[name]

	return ok
}

// RoleIDsByPrivilege - возвращает ID ролей, обладающих указанной привилегией.
func (p *PermsProvider) RoleIDsByPrivilege(name string) []uint16 {
	return p.privileges[name]
}

// HasPermission - сообщает о наличии указанного разрешения.
func (p *PermsProvider) HasPermission(name string) bool {
	_, ok := p.permissions[name]

	return ok
}

// RoleIDsByPermission - возвращает ID ролей, обладающих указанным разрешением.
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
