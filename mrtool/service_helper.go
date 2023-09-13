package mrtool

import "github.com/mondegor/go-webcore/mrcore"

type (
    ServiceHelper struct {}
)

func NewServiceHelper() *ServiceHelper {
    return &ServiceHelper{}
}

func (h *ServiceHelper) WrapErrorForSelect(err error, entityName string) error {
    if mrcore.FactoryErrStorageNoRowFound.Is(err) {
        return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Caller(1).Wrap(err, entityName)
}

func (h *ServiceHelper) WrapErrorForUpdate(err error, entityName string) error {
    if mrcore.FactoryErrStorageRowsNotAffected.Is(err) {
        return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return mrcore.FactoryErrServiceEntityNotUpdated.Caller(1).Wrap(err, entityName)
}

func (h *ServiceHelper) WrapErrorForRemove(err error, entityName string) error {
    if mrcore.FactoryErrStorageRowsNotAffected.Is(err) {
        return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return mrcore.FactoryErrServiceEntityNotRemoved.Caller(1).Wrap(err, entityName)
}

func (h *ServiceHelper) ReturnErrorIfItemNotFound(err error, entityName string) error {
    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, entityName)
        }

        return mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Caller(1).Wrap(err, entityName)
    }

    return nil
}
