package role

import "github.com/mondegor/go-webcore/mraccess"

type (
	// Role - comment struct.
	Role struct {
		roleID uint16
		rights mraccess.RightsSource
	}
)

// New - создаёт объект Role.
func New(role string, rights mraccess.RightsSource) *Role {
	if roleIDs := rights.RoleIDsByNames([]string{role}); len(roleIDs) > 0 {
		return &Role{
			roleID: roleIDs[0],
			rights: rights,
		}
	}

	return &Role{
		rights: rights,
	}
}

// CheckPrivilege - comment method.
func (r *Role) CheckPrivilege(name string) bool {
	return r.roleInArray(
		r.rights.RoleIDsByPrivilege(name),
	)
}

// CheckPermission - comment method.
func (r *Role) CheckPermission(name string) bool {
	return r.roleInArray(
		r.rights.RoleIDsByPermission(name),
	)
}

func (r *Role) roleInArray(roleIDs []uint16) bool {
	for i := range roleIDs {
		if roleIDs[i] == r.roleID {
			return true
		}
	}

	return false
}
