package mrcore

import e "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrUseCaseOperationFailed = e.NewFactory(
		"errUseCaseOperationFailed", e.ErrorTypeInternal, "operation failed")

	FactoryErrUseCaseTemporarilyUnavailable = e.NewFactory(
		"errUseCaseTemporarilyUnavailable", e.ErrorTypeSystem, "system is temporarily unavailable")

	FactoryErrUseCaseIncorrectInputData = e.NewFactory(
		"errUseCaseIncorrectInputData", e.ErrorTypeInternalNotice, "{{ .key }}={{ .data }} is incorrect")

	FactoryErrUseCaseEntityNotFound = e.NewFactory(
		"errUseCaseEntityNotFound", e.ErrorTypeUser, "entity not found")

	FactoryErrUseCaseEntityNotAvailable = e.NewFactory(
		"errUseCaseEntityNotAvailable", e.ErrorTypeUser, "entity is not available")

	FactoryErrUseCaseEntityVersionInvalid = e.NewFactory(
		"errUseCaseEntityVersionInvalid", e.ErrorTypeUser, "entity version is invalid")

	FactoryErrUseCaseSwitchStatusRejected = e.NewFactory(
		"errUseCaseSwitchStatusRejected", e.ErrorTypeUser, "switching from '{{ .statusFrom }}' to '{{ .statusTo }}' is rejected")

	FactoryErrUseCaseInvalidFile = e.NewFactory(
		"errUseCaseInvalidFile", e.ErrorTypeUser, "file is invalid")
)
