package role

import "github.com/mondegor/go-webcore/mraccess"

type (
	// Group - comment struct.
	Group struct {
		roleIDs []uint16
		rights  mraccess.RightsSource
	}
)

// NewGroup - создаёт объект Group.
func NewGroup(roles []string, rights mraccess.RightsSource) *Group {
	return &Group{
		roleIDs: rights.RoleIDsByNames(roles),
		rights:  rights,
	}
}

// CheckPrivilege - comment method.
func (g *Group) CheckPrivilege(name string) bool {
	return g.isIntersection(
		g.rights.RoleIDsByPrivilege(name),
	)
}

// CheckPermission - comment method.
func (g *Group) CheckPermission(name string) bool {
	return g.isIntersection(
		g.rights.RoleIDsByPermission(name),
	)
}

func (g *Group) isIntersection(roleIDs []uint16) bool {
	for i := range g.roleIDs {
		for j := range roleIDs {
			if g.roleIDs[i] == roleIDs[j] {
				return true
			}
		}
	}

	return false
}
