package main

import (
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
)

func main() {
	logger := mrcore.DefaultLogger()

	logger.Err(mrcore.FactoryErrInternal.New())
	logger.Err(mrcore.FactoryErrInternalTypeAssertion.New("MY-TYPE", "MY-VALUE"))

	logger.Info(mrcore.FactoryErrInternal.WithAttr("MY-DATA-KEY", mrmsg.Data{"itemId": "id-001"}).New().Error())
}
