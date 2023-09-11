package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrConnectionIsAlreadyCreated = NewFactory(
        "errStorageConnectionIsAlreadyCreated", ErrorKindInternal, "connection '{{ .name }}' is already created")

    FactoryErrConnectionIsNotOpened = NewFactory(
        "errStorageConnectionIsNotOpened", ErrorKindInternal, "connection '{{ .name }}' is not opened")

    FactoryErrConnectionFailed = NewFactory(
        "errStorageConnectionFailed", ErrorKindSystem, "connection '{{ .name }}' is failed")

    FactoryErrQueryFailed = NewFactory(
        "errStorageQueryFailed", ErrorKindInternal, "query is failed")

    FactoryErrFetchDataFailed = NewFactory(
        "errStorageFetchDataFailed", ErrorKindInternal, "fetching data is failed")

    FactoryErrFetchedInvalidData = NewFactory(
        "errStorageFetchedInvalidData", ErrorKindInternal, "fetched data '{{ .value }}' is invalid")

    FactoryErrNoRowFound = NewFactory(
        "errStorageNoRowFound", ErrorKindInternalNotice, "no row found")

    FactoryErrRowsNotAffected = NewFactory(
        "errStorageRowsNotAffected", ErrorKindInternalNotice, "rows not affected")
)
