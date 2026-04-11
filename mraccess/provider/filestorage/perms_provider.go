package filestorage

import (
	"errors"
	"fmt"
	"math"
	"slices"
)

type (
	// PermsProvider - управляет ролями, привилегиями и разрешениями подгружаемые из файлового хранилища.
	// Загружает конфигурацию ролей из YAML-файлов и предоставляет методы для проверки прав доступа.
	PermsProvider struct {
		roles       map[string]uint16   // map role names to role IDs
		privileges  map[string][]uint16 // map privilege names to role IDs
		permissions map[string][]uint16 // map permission names to role IDs
	}
)

// NewPermsProvider - создаёт объект PermsProvider, загружающий роли из указанного каталога.
// Каждая роль хранится в отдельном YAML-файле с именем роли и расширением .yaml.
// Параметры:
//   - roles - список имён ролей;
//   - opts - позволяет загружать только разрешенные привилегии и разрешения.
func NewPermsProvider(dirPath string, roles []string, opts ...Option) (*PermsProvider, error) {
	if len(roles) == 0 {
		return nil, errors.New("roles is required")
	}

	if len(roles) > math.MaxUint16 {
		return nil, errors.New("the number of roles cannot be more than 65536")
	}

	p := &PermsProvider{
		roles:       make(map[string]uint16, len(roles)),
		privileges:  make(map[string][]uint16),
		permissions: make(map[string][]uint16),
	}

	o := options{}

	for _, opt := range opts {
		opt(&o)
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
			if len(o.allowedPrivileges) > 0 && !slices.Contains(o.allowedPrivileges, priv) {
				return nil, fmt.Errorf("privilege is not registered in AllowedPrivileges (name='%s')", priv)
			}

			p.privileges[priv] = append(p.privileges[priv], roleID)
		}

		for _, perm := range fileCfg.Permissions {
			if len(o.allowedPermissions) > 0 && !slices.Contains(o.allowedPermissions, perm) {
				return nil, fmt.Errorf("permission is not registered in AllowedPermissions (name='%s')", perm)
			}

			p.permissions[perm] = append(p.permissions[perm], roleID)
		}
	}

	return p, nil
}

// RoleIDsByNames - возвращает ID ролей по их именам.
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

// HasPrivilege - сообщает, зарегистрирована ли указанная привилегия в системе.
func (p *PermsProvider) HasPrivilege(name string) bool {
	_, ok := p.privileges[name]

	return ok
}

// RoleIDsByPrivilege - возвращает ID всех ролей, обладающих указанной привилегией.
func (p *PermsProvider) RoleIDsByPrivilege(name string) []uint16 {
	return p.privileges[name]
}

// HasPermission - сообщает, зарегистрировано ли указанное разрешение в системе.
func (p *PermsProvider) HasPermission(name string) bool {
	_, ok := p.permissions[name]

	return ok
}

// RoleIDsByPermission - возвращает ID всех ролей, обладающих указанным разрешением.
func (p *PermsProvider) RoleIDsByPermission(name string) []uint16 {
	return p.permissions[name]
}
