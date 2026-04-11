package mraccess

type (
	// rolesGroupSet - внутренняя реализация набора групп ролей с привязанными к ним привилегиями и разрешениями.
	rolesGroupSet struct {
		group2rights map[string]RightsChecker // group name -> rights
	}
)

// NewRolesGroupSet - создаёт объект RightsGetter для набора групп ролей.
// Для каждой группы создаёт соответствующий RightsChecker (одну роль или группу ролей).
func NewRolesGroupSet(groups []RoleGroup, rightsSource RightsSource) (RightsGetter, error) {
	group2rights := make(map[string]RightsChecker, len(groups))

	for _, group := range groups {
		switch len(group.Roles) {
		case 0:
			// skip
		case 1:
			role, err := NewRole(group.Roles[0], rightsSource)
			if err != nil {
				return nil, err
			}

			group2rights[group.Name] = role
		default:
			roleGroup, err := NewRoleGroup(group.Roles, rightsSource)
			if err != nil {
				return nil, err
			}

			group2rights[group.Name] = roleGroup
		}
	}

	return &rolesGroupSet{
		group2rights: group2rights,
	}, nil
}

// Rights - выдаёт привилегии и разрешения для указанной группы ролей.
// Если группа не найдена, возвращает объект accessDenied.
func (g *rolesGroupSet) Rights(group string) RightsChecker {
	if rights, ok := g.group2rights[group]; ok {
		return rights
	}

	return accessDenied{}
}
