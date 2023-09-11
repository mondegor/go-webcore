package mrcore

type (
    ServiceHelper struct {}
)

func NewServiceHelper() *ServiceHelper {
    return &ServiceHelper{}
}

func (h *ServiceHelper) WrapErrorForSelect(err error, entityName string) error {
    if FactoryErrNoRowFound.Is(err) {
        return FactoryErrEntityNotFound.Wrap(err, entityName)
    }

    return FactoryErrEntityTemporarilyUnavailable.Caller(1).Wrap(err, entityName)
}

func (h *ServiceHelper) WrapErrorForUpdate(err error, entityName string) error {
    if FactoryErrRowsNotAffected.Is(err) {
        return FactoryErrEntityNotFound.Wrap(err, entityName)
    }

    return FactoryErrEntityNotUpdated.Caller(1).Wrap(err, entityName)
}

func (h *ServiceHelper) WrapErrorForRemove(err error, entityName string) error {
    if FactoryErrRowsNotAffected.Is(err) {
        return FactoryErrEntityNotFound.Wrap(err, entityName)
    }

    return FactoryErrEntityNotRemoved.Caller(1).Wrap(err, entityName)
}

func (h *ServiceHelper) ReturnErrorIfItemNotFound(err error, entityName string) error {
    if err != nil {
        if FactoryErrNoRowFound.Is(err) {
            return FactoryErrEntityNotFound.Wrap(err, entityName)
        }

        return FactoryErrEntityTemporarilyUnavailable.Caller(1).Wrap(err, entityName)
    }

    return nil
}
