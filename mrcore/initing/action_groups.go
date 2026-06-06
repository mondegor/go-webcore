package initing

import (
	"github.com/mondegor/go-sysmess/mraccess"
	"github.com/mondegor/go-sysmess/mraccess/config"
	"github.com/mondegor/go-sysmess/mrlog"
)

// InitActionGroups - создаёт и инициализирует группу действий (ActionGroup) из конфигурации.
func InitActionGroups(logger mrlog.Logger, groups []config.ActionGroup) (name2group map[string]*mraccess.ActionGroup) {
	name2group = make(map[string]*mraccess.ActionGroup, len(groups))

	for _, group := range groups {
		mrlog.Info(
			logger,
			"Init actionGroups with privilege and base path",
			"name", group.Name,
			"privilege", group.Privilege,
			"basePath", group.BasePath,
		)

		name2group[group.Name] = mraccess.NewActionGroup(
			group.Name,
			group.Privilege,
			group.BasePath,
		)
	}

	return name2group
}
