package filestorage

type (
	// PermsProviderInfo - информация о зарегистрированных ролях, привилегиях и разрешениях.
	// Используется для отладки и отображения в консоли.
	PermsProviderInfo struct {
		Roles       []string
		Privileges  []string
		Permissions []string
	}
)

// ExtractProviderInfo - извлекает данные PermsProviderInfo из PermsProvider.
func ExtractProviderInfo(provider *PermsProvider) PermsProviderInfo {
	if provider == nil {
		return PermsProviderInfo{}
	}

	return PermsProviderInfo{
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
