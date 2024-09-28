package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrStorageConnectionIsAlreadyCreated - connection is already created.
	ErrStorageConnectionIsAlreadyCreated = mrerr.NewProto(
		"errStorageConnectionIsAlreadyCreated", mrerr.ErrorKindInternal, "connection '{{ .name }}' is already created")

	// ErrStorageConnectionIsNotOpened - connection is not opened.
	ErrStorageConnectionIsNotOpened = mrerr.NewProto(
		"errStorageConnectionIsNotOpened", mrerr.ErrorKindInternal, "connection '{{ .name }}' is not opened")

	// ErrStorageConnectionIsBusy - connection is busy.
	ErrStorageConnectionIsBusy = mrerr.NewProto(
		"errStorageConnectionIsBusy", mrerr.ErrorKindInternal, "connection '{{ .name }}' is busy")

	// ErrStorageConnectionFailed - connection is failed.
	ErrStorageConnectionFailed = mrerr.NewProto(
		"errStorageConnectionFailed", mrerr.ErrorKindSystem, "connection '{{ .name }}' is failed")

	// ErrStorageQueryFailed - query is failed.
	ErrStorageQueryFailed = mrerr.NewProto(
		"errStorageQueryFailed", mrerr.ErrorKindInternal, "query is failed")

	// ErrStorageFetchDataFailed - fetching data is failed.
	ErrStorageFetchDataFailed = mrerr.NewProto(
		"errStorageFetchDataFailed", mrerr.ErrorKindInternal, "fetching data is failed")

	// ErrStorageNoRowFound - no row found.
	ErrStorageNoRowFound = mrerr.NewProto(
		"errStorageNoRowFound", mrerr.ErrorKindInternal, "no row found")

	// ErrStorageRowsNotAffected - rows not affected.
	ErrStorageRowsNotAffected = mrerr.NewProto(
		"errStorageRowsNotAffected", mrerr.ErrorKindInternal, "rows not affected")
)
