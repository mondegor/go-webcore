package main

import (
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

func main() {
	logger := mrlog.New(mrlog.DebugLevel).With().Str("example", "errors").Logger()

	logger.Error().Err(mrcore.FactoryErrInternal.New()).Msg("this is FactoryErrInternal")
	logger.Error().Err(mrcore.FactoryErrInternalTypeAssertion.New("MY-TYPE", "MY-VALUE")).Send()
	logger.Error().Err(mrcore.FactoryErrInternal.WithAttr("MY-DATA-KEY", mrmsg.Data{"itemId": "id-001"}).New()).Send()

	logger.Fatal().Int("int1", 1).Int("int2", 2).Int("int3", 3).Msg(mrcore.FactoryErrInternal.New().Error())
}
