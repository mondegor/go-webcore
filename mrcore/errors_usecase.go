package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrUseCaseOperationFailed = NewFactoryWithCaller(
		"errUseCaseOperationFailed", ErrorKindInternal, "operation failed")

	FactoryErrUseCaseTemporarilyUnavailable = NewFactoryWithCaller(
		"errUseCaseTemporarilyUnavailable", ErrorKindSystem, "system is temporarily unavailable")

	FactoryErrUseCaseIncorrectInputData = NewFactory(
		"errUseCaseIncorrectInputData", ErrorKindInternal, "{{ .key }}={{ .data }} is incorrect")

	FactoryErrUseCaseEntityNotFound = NewFactory(
		"errUseCaseEntityNotFound", ErrorKindUser, "entity not found")

	FactoryErrUseCaseEntityNotAvailable = NewFactory(
		"errUseCaseEntityNotAvailable", ErrorKindUser, "entity is not available")

	FactoryErrUseCaseEntityVersionInvalid = NewFactory(
		"errUseCaseEntityVersionInvalid", ErrorKindUser, "entity version is invalid")

	FactoryErrUseCaseSwitchStatusRejected = NewFactory(
		"errUseCaseSwitchStatusRejected", ErrorKindUser, "switching from '{{ .statusFrom }}' to '{{ .statusTo }}' is rejected")

	FactoryErrUseCaseInvalidFile = NewFactory(
		"errUseCaseInvalidFile", ErrorKindUser, "file is invalid")
)
