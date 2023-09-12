package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrServiceIncorrectInputData = NewFactory(
        "errServiceIncorrectData", ErrorKindInternal, "data '{{ .data }}' is incorrect")

    FactoryErrServiceEntityTemporarilyUnavailable = NewFactory(
        "errServiceEntityTemporarilyUnavailable", ErrorKindSystem, "entity '{{ .name }}' is temporarily unavailable")

    FactoryErrServiceEntityNotFound = NewFactory(
        "errServiceEntityNotFound", ErrorKindInternalNotice, "entity '{{ .name }}' is not found")

    FactoryErrServiceEntityNotCreated = NewFactory(
        "errServiceEntityNotCreated", ErrorKindSystem, "entity '{{ .name }}' is not created")

    FactoryErrServiceEntityNotUpdated = NewFactory(
        "errServiceEntityNotUpdated", ErrorKindSystem, "entity '{{ .name }}' is not updated")

    FactoryErrServiceEntityNotRemoved = NewFactory(
        "errServiceEntityNotRemoved", ErrorKindSystem, "entity '{{ .name }}' is not removed")

    FactoryErrServiceIncorrectSwitchStatus = NewFactory(
        "errServiceIncorrectSwitchStatus", ErrorKindInternal, "incorrect switch status: '{{ .currentStatus }}' -> '{{ .statusTo }}' for entity '{{ .name }}(ID={{ .id }})'")
)
