package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrServiceEmptyInputData = NewFactory(
		"errServiceEmptyData", ErrorKindInternalNotice, "value of '{{ .name }}' is empty")

	FactoryErrServiceIncorrectInputData = NewFactory(
		"errServiceIncorrectData", ErrorKindInternalNotice, "data '{{ .data }}' is incorrect")

	FactoryErrServiceTemporarilyUnavailable = NewFactory(
		"errServiceTemporarilyUnavailable", ErrorKindSystem, "resource '{{ .data }}' is temporarily unavailable")

	FactoryErrServiceEntityNotFound = NewFactory(
		"errServiceEntityNotFound", ErrorKindInternalNotice, "entity '{{ .data }}' is not found")

	FactoryErrServiceEntityNotCreated = NewFactory(
		"errServiceEntityNotCreated", ErrorKindInternal, "entity '{{ .data }}' is not created")

	FactoryErrServiceEntityNotStored = NewFactory(
		"errServiceEntityNotStored", ErrorKindInternal, "entity '{{ .data }}' is not stored")

	FactoryErrServiceEntityNotRemoved = NewFactory(
		"errServiceEntityNotRemoved", ErrorKindInternal, "entity '{{ .data }}' is not removed")

	FactoryErrServiceEntitySwitchStatusImpossible = NewFactory(
		"errServiceEntitySwitchStatusImpossible", ErrorKindInternalNotice, "entity '{{ .data }}': switching from '{{ .statusFrom }}' to '{{ .statusTo }}' is impossible")
)
