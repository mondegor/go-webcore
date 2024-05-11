package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrStorageConnectionIsAlreadyCreated = NewFactory(
		"errStorageConnectionIsAlreadyCreated", ErrorTypeInternal, "connection '{{ .name }}' is already created")

	FactoryErrStorageConnectionIsNotOpened = NewFactory(
		"errStorageConnectionIsNotOpened", ErrorTypeInternal, "connection '{{ .name }}' is not opened")

	FactoryErrStorageConnectionFailed = NewFactory(
		"errStorageConnectionFailed", ErrorTypeSystem, "connection '{{ .name }}' is failed")

	FactoryErrStorageQueryFailed = NewFactory(
		"errStorageQueryFailed", ErrorTypeInternal, "query is failed")

	FactoryErrStorageFetchDataFailed = NewFactory(
		"errStorageFetchDataFailed", ErrorTypeInternal, "fetching data is failed")

	FactoryErrStorageNoRowFound = NewFactory(
		"errStorageNoRowFound", ErrorTypeInternalNotice, "no row found")

	FactoryErrStorageRowsNotAffected = NewFactory(
		"errStorageRowsNotAffected", ErrorTypeInternalNotice, "rows not affected")
)
