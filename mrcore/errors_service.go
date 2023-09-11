package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrIncorrectInputData = NewFactory(
        "errServiceIncorrectData", ErrorKindInternal, "data '{{ .data }}' is incorrect")

    FactoryErrEntityTemporarilyUnavailable = NewFactory(
        "errServiceEntityTemporarilyUnavailable", ErrorKindSystem, "entity '{{ .name }}' is temporarily unavailable")

    FactoryErrEntityNotFound = NewFactory(
        "errServiceEntityNotFound", ErrorKindInternalNotice, "entity '{{ .name }}' is not found")

    FactoryErrEntityNotCreated = NewFactory(
        "errServiceEntityNotCreated", ErrorKindSystem, "entity '{{ .name }}' is not created")

    FactoryErrEntityNotUpdated = NewFactory(
        "errServiceEntityNotUpdated", ErrorKindSystem, "entity '{{ .name }}' is not updated")

    FactoryErrEntityNotRemoved = NewFactory(
        "errServiceEntityNotRemoved", ErrorKindSystem, "entity '{{ .name }}' is not removed")

    FactoryErrIncorrectSwitchStatus = NewFactory(
        "errServiceIncorrectSwitchStatus", ErrorKindInternal, "incorrect switch status: '{{ .currentStatus }}' -> '{{ .statusTo }}' for entity '{{ .name }}(ID={{ .id }})'")
)
