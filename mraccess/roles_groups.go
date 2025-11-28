package mraccess

type (
	// rolesGroupSet - список групп ролей с привязанными к ним привилегиями и разрешениями.
	rolesGroupSet struct {
		group2rights map[string]RightsChecker // group name -> rights
	}
)

// NewRolesGroupSet - создаёт объект RightsGetter.
func NewRolesGroupSet(groups []RoleGroup, rightsSource RightsSource) RightsGetter {
	group2rights := make(map[string]RightsChecker, len(groups))

	for _, group := range groups {
		if len(group.Roles) == 0 {
			continue
		}

		if len(group.Roles) == 1 {
			group2rights[group.Name] = NewRole(group.Roles[0], rightsSource)
		} else {
			group2rights[group.Name] = NewRoleGroup(group.Roles, rightsSource)
		}
	}

	return &rolesGroupSet{
		group2rights: group2rights,
	}
}

// Rights - выдаёт привилегии и разрешения для указанной группы ролей.
func (g *rolesGroupSet) Rights(group string) RightsChecker {
	if rights, ok := g.group2rights[group]; ok {
		return rights
	}

	return accessDenied{}
}
