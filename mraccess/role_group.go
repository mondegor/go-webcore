package mraccess

import "fmt"

type (
	// RoleGroup - группа с привязанными к ней ролями.
	// Используется для объединения нескольких ролей под одним именем группы.
	RoleGroup struct {
		Name  string
		Roles []string
	}

	// entryRoleGroup - внутренняя реализация группы ролей с привязанными
	// к ним привилегиями и разрешениями.
	entryRoleGroup struct {
		roleIDs []uint16
		rights  RightsSource
	}
)

// NewRoleGroup - создаёт объект RightsChecker для группы ролей.
func NewRoleGroup(roles []string, rights RightsSource) (RightsChecker, error) {
	roleIDs := rights.RoleIDsByNames(roles)

	if len(roleIDs) != len(roles) {
		return nil, fmt.Errorf("there are fewer role IDs than role names: roles=%v", roles)
	}

	return &entryRoleGroup{
		roleIDs: roleIDs,
		rights:  rights,
	}, nil
}

// HasPrivilege - сообщает, обладает хотя бы одна роль в группе указанной привилегией.
func (g *entryRoleGroup) HasPrivilege(name string) bool {
	return g.isIntersection(
		g.rights.RoleIDsByPrivilege(name),
	)
}

// HasPermission - сообщает, обладает хотя бы одна роль в группе указанным разрешением.
func (g *entryRoleGroup) HasPermission(name string) bool {
	return g.isIntersection(
		g.rights.RoleIDsByPermission(name),
	)
}

func (g *entryRoleGroup) isIntersection(roleIDs []uint16) bool {
	for i := range g.roleIDs {
		for j := range roleIDs {
			if g.roleIDs[i] == roleIDs[j] {
				return true
			}
		}
	}

	return false
}
