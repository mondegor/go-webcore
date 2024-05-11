package mrcore

import e "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrStorageConnectionIsAlreadyCreated = e.NewFactory(
		"errStorageConnectionIsAlreadyCreated", e.ErrorTypeInternal, "connection '{{ .name }}' is already created")

	FactoryErrStorageConnectionIsNotOpened = e.NewFactory(
		"errStorageConnectionIsNotOpened", e.ErrorTypeInternal, "connection '{{ .name }}' is not opened")

	FactoryErrStorageConnectionFailed = e.NewFactory(
		"errStorageConnectionFailed", e.ErrorTypeSystem, "connection '{{ .name }}' is failed")

	FactoryErrStorageQueryFailed = e.NewFactory(
		"errStorageQueryFailed", e.ErrorTypeInternal, "query is failed")

	FactoryErrStorageFetchDataFailed = e.NewFactory(
		"errStorageFetchDataFailed", e.ErrorTypeInternal, "fetching data is failed")

	FactoryErrStorageNoRowFound = e.NewFactory(
		"errStorageNoRowFound", e.ErrorTypeInternalNotice, "no row found")

	FactoryErrStorageRowsNotAffected = e.NewFactory(
		"errStorageRowsNotAffected", e.ErrorTypeInternalNotice, "rows not affected")
)
