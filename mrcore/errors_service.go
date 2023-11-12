package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrServiceEmptyInputData = NewFactory( // :TODO: wrapErrorFunc
		"errServiceEmptyData", ErrorKindInternalNotice, "{{ .name }} is empty")

	FactoryErrServiceIncorrectInputData = NewFactory( // :TODO: wrapErrorFunc
		"errServiceIncorrectData", ErrorKindInternalNotice, "data '{{ .data }}' is incorrect")

	FactoryErrServiceTemporarilyUnavailable = NewFactory(
		"errServiceTemporarilyUnavailable", ErrorKindSystem, "resource '{{ .name }}' is temporarily unavailable")

	FactoryErrServiceEntityNotFound = NewFactory(
		"errServiceEntityNotFound", ErrorKindInternalNotice, "entity '{{ .name }}' is not found")

	FactoryErrServiceEntityNotCreated = NewFactory(
		"errServiceEntityNotCreated", ErrorKindSystem, "entity '{{ .name }}' is not created")

	FactoryErrServiceEntityNotUpdated = NewFactory(
		"errServiceEntityNotUpdated", ErrorKindSystem, "entity '{{ .name }}' is not updated")

	FactoryErrServiceEntityNotRemoved = NewFactory(
		"errServiceEntityNotRemoved", ErrorKindSystem, "entity '{{ .name }}' is not removed")

	FactoryErrServiceIncorrectSwitchStatus = NewFactory(
		"errServiceIncorrectSwitchStatus", ErrorKindInternalNotice, "incorrect switch status: '{{ .currentStatus }}' -> '{{ .statusTo }}' for entity '{{ .name }}(ID={{ .id }})'")
)
