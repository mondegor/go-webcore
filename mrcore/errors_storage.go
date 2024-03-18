package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrStorageConnectionIsAlreadyCreated = NewFactoryWithCaller(
		"errStorageConnectionIsAlreadyCreated", ErrorKindInternal, "connection '{{ .name }}' is already created")

	FactoryErrStorageConnectionIsNotOpened = NewFactoryWithCaller(
		"errStorageConnectionIsNotOpened", ErrorKindInternal, "connection '{{ .name }}' is not opened")

	FactoryErrStorageConnectionFailed = NewFactoryWithCaller(
		"errStorageConnectionFailed", ErrorKindSystem, "connection '{{ .name }}' is failed")

	FactoryErrStorageQueryFailed = NewFactoryWithCaller(
		"errStorageQueryFailed", ErrorKindInternal, "query is failed")

	FactoryErrStorageFetchDataFailed = NewFactoryWithCaller(
		"errStorageFetchDataFailed", ErrorKindInternal, "fetching data is failed")

	FactoryErrStorageNoRowFound = NewFactory(
		"errStorageNoRowFound", ErrorKindInternal, "no row found")

	FactoryErrStorageRowsNotAffected = NewFactory(
		"errStorageRowsNotAffected", ErrorKindInternal, "rows not affected")
)
