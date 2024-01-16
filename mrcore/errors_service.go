package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrServiceOperationFailed = NewFactory(
		"errServiceOperationFailed", ErrorKindInternal, "operation failed")

	FactoryErrServiceTemporarilyUnavailable = NewFactory(
		"errServiceTemporarilyUnavailable", ErrorKindSystem, "system is temporarily unavailable")

	FactoryErrServiceIncorrectInputData = NewFactory(
		"errServiceIncorrectInputData", ErrorKindInternalNotice, "{{ .key }}={{ .data }} is incorrect")

	FactoryErrServiceEntityNotFound = NewFactory(
		"errServiceEntityNotFound", ErrorKindUser, "entity not found")

	FactoryErrServiceEntityVersionInvalid = NewFactory(
		"errServiceEntityVersionInvalid", ErrorKindUser, "entity version is invalid")

	FactoryErrServiceSwitchStatusRejected = NewFactory(
		"errServiceSwitchStatusRejected", ErrorKindUser, "switching from '{{ .statusFrom }}' to '{{ .statusTo }}' is rejected")

	FactoryErrServiceInvalidFile = NewFactory(
		"errServiceInvalidFile", ErrorKindUser, "file is invalid")
)
