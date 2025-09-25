package mrinit

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
)

// InfoCreateModule - comment func.
func InfoCreateModule(logger mrlog.LiteLogger, name string) {
	logger.Info(fmt.Sprintf("Create and init module '%s'", name))
}

// InfoCreateUnit - comment func.
func InfoCreateUnit(logger mrlog.LiteLogger, name string) {
	logger.Info(fmt.Sprintf("Create and init unit '%s' of the module", name))
}
