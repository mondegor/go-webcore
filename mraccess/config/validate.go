package config

import (
	"fmt"
	"slices"
)

const (
	// PublicPrivilege - привилегия публичного доступа.
	PublicPrivilege = "public"
)

// ValidateActionGroups - валидация групп обработчиков.
func ValidateActionGroups(actionGroups []ActionGroup, allPrivileges []string) error {
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

		if group.Privilege != PublicPrivilege {
			if !slices.Contains(allPrivileges, group.Privilege) {
				return fmt.Errorf("privilege is not found in privileges for actionGroup (name='%s', group='%s')", group.Privilege, group.Name)
			}
		}
	}

	return nil
}
