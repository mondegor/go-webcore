package groups

import (
	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/role"
)

type (
	// Groups - comment struct.
	Groups struct {
		rightsMap map[string]mraccess.RightsChecker // group name -> rights
	}

	// Group - comment struct.
	Group struct {
		Name  string
		Roles []string
	}
)

// NewGroups - создаёт объект Groups.
func NewGroups(groups []Group, rightsSource mraccess.RightsSource) *Groups {
	rightsMap := make(map[string]mraccess.RightsChecker, len(groups))

	for _, group := range groups {
		if len(group.Roles) == 0 {
			continue
		}

		if len(group.Roles) == 1 {
			rightsMap[group.Name] = role.New(group.Roles[0], rightsSource)
		} else {
			rightsMap[group.Name] = role.NewGroup(group.Roles, rightsSource)
		}
	}

	return &Groups{
		rightsMap: rightsMap,
	}
}

// Rights - comment method.
func (g *Groups) Rights(name string) mraccess.RightsChecker {
	if rights, ok := g.rightsMap[name]; ok {
		return rights
	}

	return accessDenied{}
}
