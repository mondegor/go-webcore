package config

import (
	"fmt"
	"slices"
)

const (
	// PublicPrivilege - привилегия публичного доступа, не требующая аутентификации.
	PublicPrivilege = "public"
)

// ValidateActionGroups - выполняет валидацию конфигурации групп обработчиков.
// Проверяет уникальность имён и путей групп, а также наличие привилегий групп обработчиков
// в предоставленном списке allPrivileges.
func ValidateActionGroups(actionGroups []ActionGroup, allowedPrivileges []string) error {
	uniqNames := make(map[string]bool, len(actionGroups))
	uniqPaths := make(map[string]bool, len(actionGroups))

	for _, group := range actionGroups {
		if uniqNames[group.Name] {
			return fmt.Errorf("duplicate actionGroup name '%s'", group.Name)
		}

		if uniqPaths[group.BasePath] {
			return fmt.Errorf("duplicate base path for actionGroup (path='%s', group='%s')", group.BasePath, group.Name)
		}

		uniqNames[group.Name] = true
		uniqPaths[group.BasePath] = true

		if group.Privilege == PublicPrivilege {
			continue
		}

		if !slices.Contains(allowedPrivileges, group.Privilege) {
			return fmt.Errorf("privilege is not found in privileges for actionGroup (name='%s', group='%s')", group.Privilege, group.Name)
		}
	}

	return nil
}
