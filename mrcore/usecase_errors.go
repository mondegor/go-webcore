package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrUseCaseOperationFailed - operation failed.
	ErrUseCaseOperationFailed = mrerr.NewProto(
		"errUseCaseOperationFailed", mrerr.ErrorKindInternal, "operation failed")

	// ErrUseCaseTemporarilyUnavailable - system is temporarily unavailable.
	ErrUseCaseTemporarilyUnavailable = mrerr.NewProto(
		"errUseCaseTemporarilyUnavailable", mrerr.ErrorKindSystem, "system is temporarily unavailable")

	// ErrUseCaseIncorrectInputData - data is incorrect.
	ErrUseCaseIncorrectInputData = mrerr.NewProto(
		"errUseCaseIncorrectInputData", mrerr.ErrorKindInternal, "{{ .key }}={{ .data }} is incorrect")

	// ErrUseCaseEntityNotFound - entity not found.
	ErrUseCaseEntityNotFound = mrerr.NewProto(
		"errUseCaseEntityNotFound", mrerr.ErrorKindUser, "entity not found")

	// ErrUseCaseEntityNotAvailable - entity is not available.
	ErrUseCaseEntityNotAvailable = mrerr.NewProto(
		"errUseCaseEntityNotAvailable", mrerr.ErrorKindUser, "entity is not available")

	// ErrUseCaseEntityVersionInvalid - entity version is invalid.
	ErrUseCaseEntityVersionInvalid = mrerr.NewProto(
		"errUseCaseEntityVersionInvalid", mrerr.ErrorKindUser, "entity version is invalid")

	// ErrUseCaseSwitchStatusRejected - switching from status to status is rejected.
	ErrUseCaseSwitchStatusRejected = mrerr.NewProto(
		"errUseCaseSwitchStatusRejected", mrerr.ErrorKindUser, "switching from '{{ .statusFrom }}' to '{{ .statusTo }}' is rejected")

	// ErrUseCaseInvalidFile - file is invalid.
	ErrUseCaseInvalidFile = mrerr.NewProto(
		"errUseCaseInvalidFile", mrerr.ErrorKindUser, "file is invalid")
)
