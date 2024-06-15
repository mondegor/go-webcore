package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrUseCaseOperationFailed - comment error.
	ErrUseCaseOperationFailed = mrerr.NewProto(
		"errUseCaseOperationFailed", mrerr.ErrorKindInternal, "operation failed")

	// ErrUseCaseTemporarilyUnavailable - comment error.
	ErrUseCaseTemporarilyUnavailable = mrerr.NewProto(
		"errUseCaseTemporarilyUnavailable", mrerr.ErrorKindSystem, "system is temporarily unavailable")

	// ErrUseCaseIncorrectInputData - comment error.
	ErrUseCaseIncorrectInputData = mrerr.NewProto(
		"errUseCaseIncorrectInputData", mrerr.ErrorKindInternal, "{{ .key }}={{ .data }} is incorrect")

	// ErrUseCaseEntityNotFound - comment error.
	ErrUseCaseEntityNotFound = mrerr.NewProto(
		"errUseCaseEntityNotFound", mrerr.ErrorKindUser, "entity not found")

	// ErrUseCaseEntityNotAvailable - comment error.
	ErrUseCaseEntityNotAvailable = mrerr.NewProto(
		"errUseCaseEntityNotAvailable", mrerr.ErrorKindUser, "entity is not available")

	// ErrUseCaseEntityVersionInvalid - comment error.
	ErrUseCaseEntityVersionInvalid = mrerr.NewProto(
		"errUseCaseEntityVersionInvalid", mrerr.ErrorKindUser, "entity version is invalid")

	// ErrUseCaseSwitchStatusRejected - comment error.
	ErrUseCaseSwitchStatusRejected = mrerr.NewProto(
		"errUseCaseSwitchStatusRejected", mrerr.ErrorKindUser, "switching from '{{ .statusFrom }}' to '{{ .statusTo }}' is rejected")

	// ErrUseCaseInvalidFile - comment error.
	ErrUseCaseInvalidFile = mrerr.NewProto(
		"errUseCaseInvalidFile", mrerr.ErrorKindUser, "file is invalid")
)
