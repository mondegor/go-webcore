package main

import (
	"github.com/mondegor/go-sysmess/mrmsg"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/mrlogbase"
)

func main() {
	logger := mrlogbase.New(mrlog.DebugLevel).With().Str("example", "errors").Logger()

	logger.Error().Err(mrcore.ErrInternal).Msg("this is ErrInternal")
	logger.Error().Err(mrcore.ErrInternalTypeAssertion.New("MY-TYPE", "MY-VALUE")).Send()
	logger.Error().Err(mrcore.ErrInternal.New().WithAttr("MY-DATA-KEY", mrmsg.Data{"itemId": "id-001"})).Send()

	logger.Fatal().Int("int1", 1).Int("int2", 2).Int("int3", 3).Msg(mrcore.ErrInternal.New().Error())
}
