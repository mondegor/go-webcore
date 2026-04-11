package mraccess

import (
	"fmt"
)

type (
	// entryRole - внутренняя реализация роли с привязанными
	// к ней привилегиями и разрешениями.
	entryRole struct {
		roleID uint16
		rights RightsSource
	}
)

// NewRole - создаёт объект RightsChecker для одной роли.
func NewRole(role string, rights RightsSource) (RightsChecker, error) {
	roleIDs := rights.RoleIDsByNames([]string{role})
	if len(roleIDs) == 0 {
		return nil, fmt.Errorf("no role found: name=%s", role)
	}

	return &entryRole{
		roleID: roleIDs[0],
		rights: rights,
	}, nil
}

// HasPrivilege - сообщает, обладает ли роль указанной привилегией.
func (r *entryRole) HasPrivilege(name string) bool {
	return r.roleInArray(
		r.rights.RoleIDsByPrivilege(name),
	)
}

// HasPermission - сообщает, обладает ли роль указанным разрешением.
func (r *entryRole) HasPermission(name string) bool {
	return r.roleInArray(
		r.rights.RoleIDsByPermission(name),
	)
}

func (r *entryRole) roleInArray(roleIDs []uint16) bool {
	// TODO: можно будет сделать бинарный поиск,
	//  если ролей в проекте достаточно большое кол-во
	for i := range roleIDs {
		if roleIDs[i] == r.roleID {
			return true
		}
	}

	return false
}
