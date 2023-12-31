package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrStorageConnectionIsAlreadyCreated = NewFactory(
		"errStorageConnectionIsAlreadyCreated", ErrorKindInternal, "connection '{{ .name }}' is already created")

	FactoryErrStorageConnectionIsNotOpened = NewFactory(
		"errStorageConnectionIsNotOpened", ErrorKindInternal, "connection '{{ .name }}' is not opened")

	FactoryErrStorageConnectionFailed = NewFactory(
		"errStorageConnectionFailed", ErrorKindSystem, "connection '{{ .name }}' is failed")

	FactoryErrStorageQueryFailed = NewFactory(
		"errStorageQueryFailed", ErrorKindInternal, "query is failed")

	FactoryErrStorageFetchDataFailed = NewFactory(
		"errStorageFetchDataFailed", ErrorKindInternal, "fetching data is failed")

	FactoryErrStorageNoRowFound = NewFactory(
		"errStorageNoRowFound", ErrorKindInternalNotice, "no row found")

	FactoryErrStorageRowsNotAffected = NewFactory(
		"errStorageRowsNotAffected", ErrorKindInternalNotice, "rows not affected")
)
