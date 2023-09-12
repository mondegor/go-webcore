package mrcore

type (
    ServiceHelper struct {}
)

func NewServiceHelper() *ServiceHelper {
    return &ServiceHelper{}
}

func (h *ServiceHelper) WrapErrorForSelect(err error, entityName string) error {
    if FactoryErrStorageNoRowFound.Is(err) {
        return FactoryErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return FactoryErrServiceEntityTemporarilyUnavailable.Caller(1).Wrap(err, entityName)
}

func (h *ServiceHelper) WrapErrorForUpdate(err error, entityName string) error {
    if FactoryErrStorageRowsNotAffected.Is(err) {
        return FactoryErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return FactoryErrServiceEntityNotUpdated.Caller(1).Wrap(err, entityName)
}

func (h *ServiceHelper) WrapErrorForRemove(err error, entityName string) error {
    if FactoryErrStorageRowsNotAffected.Is(err) {
        return FactoryErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return FactoryErrServiceEntityNotRemoved.Caller(1).Wrap(err, entityName)
}

func (h *ServiceHelper) ReturnErrorIfItemNotFound(err error, entityName string) error {
    if err != nil {
        if FactoryErrStorageNoRowFound.Is(err) {
            return FactoryErrServiceEntityNotFound.Wrap(err, entityName)
        }

        return FactoryErrServiceEntityTemporarilyUnavailable.Caller(1).Wrap(err, entityName)
    }

    return nil
}
