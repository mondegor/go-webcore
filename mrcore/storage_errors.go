package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrStorageConnectionIsAlreadyCreated - comment error.
	ErrStorageConnectionIsAlreadyCreated = mrerr.NewProto(
		"errStorageConnectionIsAlreadyCreated", mrerr.ErrorKindInternal, "connection '{{ .name }}' is already created")

	// ErrStorageConnectionIsNotOpened - comment error.
	ErrStorageConnectionIsNotOpened = mrerr.NewProto(
		"errStorageConnectionIsNotOpened", mrerr.ErrorKindInternal, "connection '{{ .name }}' is not opened")

	// ErrStorageConnectionFailed - comment error.
	ErrStorageConnectionFailed = mrerr.NewProto(
		"errStorageConnectionFailed", mrerr.ErrorKindSystem, "connection '{{ .name }}' is failed")

	// ErrStorageQueryFailed - comment error.
	ErrStorageQueryFailed = mrerr.NewProto(
		"errStorageQueryFailed", mrerr.ErrorKindInternal, "query is failed")

	// ErrStorageFetchDataFailed - comment error.
	ErrStorageFetchDataFailed = mrerr.NewProto(
		"errStorageFetchDataFailed", mrerr.ErrorKindInternal, "fetching data is failed")

	// ErrStorageNoRowFound - comment error.
	ErrStorageNoRowFound = mrerr.NewProto(
		"errStorageNoRowFound", mrerr.ErrorKindInternal, "no row found")

	// ErrStorageRowsNotAffected - comment error.
	ErrStorageRowsNotAffected = mrerr.NewProto(
		"errStorageRowsNotAffected", mrerr.ErrorKindInternal, "rows not affected")
)
