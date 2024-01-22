package mrfactory

import (
	"github.com/mondegor/go-webcore/mrcore"
)

func InfoCreateModule(logger mrcore.Logger, name string) {
	logger.Info("Create module '%s'", name)
}

func InfoCreateUnit(logger mrcore.Logger, name string) {
	logger.Info("Create unit '%s' of the module", name)
}
