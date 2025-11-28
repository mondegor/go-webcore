package config

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlib/extstrings"
)

const (
	// PublicPrivilege - привилегия публичного доступа.
	PublicPrivilege = "public"
)

// ValidateActionGroups - валидация групп обработчиков.
func ValidateActionGroups(actionGroups []ActionGroup, allPrivileges []string) error {
	uniqNames := make(map[string]struct{}, len(actionGroups))
	uniqPaths := make(map[string]struct{}, len(actionGroups))

	for _, group := range actionGroups {
		if _, ok := uniqNames[group.Name]; ok {
			return fmt.Errorf("duplicate actionGroup name '%s'", group.Name)
		}

		if _, ok := uniqPaths[group.BasePath]; ok {
			return fmt.Errorf("duplicate base path for actionGroup (path='%s', group='%s')", group.BasePath, group.Name)
		}

		uniqNames[group.Name] = struct{}{}
		uniqPaths[group.BasePath] = struct{}{}

		if group.Privilege != PublicPrivilege {
			if !extstrings.InArray(group.Privilege, allPrivileges) {
				return fmt.Errorf("privilege is not found in privileges for actionGroup (name='%s', group='%s')", group.Privilege, group.Name)
			}
		}
	}

	return nil
}
