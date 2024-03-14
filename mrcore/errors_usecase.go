package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrUseCaseOperationFailed = NewFactory(
		"errUseCaseOperationFailed", ErrorKindInternal, "operation failed")

	FactoryErrUseCaseTemporarilyUnavailable = NewFactory(
		"errUseCaseTemporarilyUnavailable", ErrorKindSystem, "system is temporarily unavailable")

	FactoryErrUseCaseIncorrectInputData = NewFactory(
		"errUseCaseIncorrectInputData", ErrorKindInternalNotice, "{{ .key }}={{ .data }} is incorrect")

	FactoryErrUseCaseEntityNotFound = NewFactory(
		"errUseCaseEntityNotFound", ErrorKindUser, "entity not found")

	FactoryErrUseCaseEntityVersionInvalid = NewFactory(
		"errUseCaseEntityVersionInvalid", ErrorKindUser, "entity version is invalid")

	FactoryErrUseCaseSwitchStatusRejected = NewFactory(
		"errUseCaseSwitchStatusRejected", ErrorKindUser, "switching from '{{ .statusFrom }}' to '{{ .statusTo }}' is rejected")

	FactoryErrUseCaseInvalidFile = NewFactory(
		"errUseCaseInvalidFile", ErrorKindUser, "file is invalid")
)
