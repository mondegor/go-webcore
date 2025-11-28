package filestorage

type (
	// RegisteredPermsInfo - информация о зарегистрированных ролях и разрешениях.
	RegisteredPermsInfo struct {
		Roles       []string
		Privileges  []string
		Permissions []string
	}
)

// NewRegisteredPermsInfo - создаёт объект RegisteredPermsInfo.
func NewRegisteredPermsInfo(provider *PermsProvider) RegisteredPermsInfo {
	if provider == nil {
		return RegisteredPermsInfo{}
	}

	return RegisteredPermsInfo{
		Roles:       registeredRoles(provider.roles),
		Privileges:  registeredPermissions(provider.privileges),
		Permissions: registeredPermissions(provider.permissions),
	}
}

func registeredRoles(roles map[string]uint16) []string {
	list := make([]string, 0, len(roles))

	for name := range roles {
		list = append(list, name)
	}

	return list
}

func registeredPermissions(permissions map[string][]uint16) []string {
	list := make([]string, 0, len(permissions))

	for name := range permissions {
		list = append(list, name)
	}

	return list
}
