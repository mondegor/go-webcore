package main

import "github.com/mondegor/go-webcore/mrcore"

func main() {
	logger := mrcore.DefaultLogger().DisableFileLine()

	logger.Err(mrcore.FactoryErrInternal.New())
	logger.Err(mrcore.FactoryErrInternalTypeAssertion.New("MY-TYPE", "MY-VALUE"))

	logger.Info(mrcore.FactoryErrInternalWithData.New("MY-DATA").Error())
}
