package main

import "github.com/mondegor/go-webcore/mrcore"

func main() {
	logger := mrcore.DefaultLogger().DisableFileLine()

	logger.Err(mrcore.FactoryErrInternal.Caller(-1).New())
	logger.Err(mrcore.FactoryErrInternalTypeAssertion.Caller(-1).New("MY-TYPE", "MY-VALUE"))

	logger.Info(mrcore.FactoryErrStorageQueryDataContainer.Caller(-1).New("MY-DATA").Error())
}
