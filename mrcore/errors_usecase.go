package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrUseCaseOperationFailed = NewFactory(
		"errUseCaseOperationFailed", ErrorTypeInternal, "operation failed")

	FactoryErrUseCaseTemporarilyUnavailable = NewFactory(
		"errUseCaseTemporarilyUnavailable", ErrorTypeSystem, "system is temporarily unavailable")

	FactoryErrUseCaseIncorrectInputData = NewFactory(
		"errUseCaseIncorrectInputData", ErrorTypeInternalNotice, "{{ .key }}={{ .data }} is incorrect")

	FactoryErrUseCaseEntityNotFound = NewFactory(
		"errUseCaseEntityNotFound", ErrorTypeUser, "entity not found")

	FactoryErrUseCaseEntityNotAvailable = NewFactory(
		"errUseCaseEntityNotAvailable", ErrorTypeUser, "entity is not available")

	FactoryErrUseCaseEntityVersionInvalid = NewFactory(
		"errUseCaseEntityVersionInvalid", ErrorTypeUser, "entity version is invalid")

	FactoryErrUseCaseSwitchStatusRejected = NewFactory(
		"errUseCaseSwitchStatusRejected", ErrorTypeUser, "switching from '{{ .statusFrom }}' to '{{ .statusTo }}' is rejected")

	FactoryErrUseCaseInvalidFile = NewFactory(
		"errUseCaseInvalidFile", ErrorTypeUser, "file is invalid")
)
