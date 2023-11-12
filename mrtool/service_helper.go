package mrtool

import "github.com/mondegor/go-webcore/mrcore"

type (
    ServiceHelper struct {
        callerSkip int
    }
)

func NewServiceHelper() *ServiceHelper {
    return &ServiceHelper{
        callerSkip: 1,
    }
}

func (h *ServiceHelper) Caller(skip int) *ServiceHelper {
    return &ServiceHelper{
        callerSkip: h.callerSkip + skip,
    }
}

func (h *ServiceHelper) WrapErrorForSelect(err error, entityName string) error {
    if mrcore.FactoryErrStorageNoRowFound.Is(err) {
        return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return mrcore.FactoryErrServiceTemporarilyUnavailable.Caller(h.callerSkip).Wrap(err, entityName)
}

func (h *ServiceHelper) WrapErrorForUpdate(err error, entityName string) error {
    if mrcore.FactoryErrStorageRowsNotAffected.Is(err) {
        return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return mrcore.FactoryErrServiceEntityNotUpdated.Caller(h.callerSkip).Wrap(err, entityName)
}

func (h *ServiceHelper) WrapErrorForRemove(err error, entityName string) error {
    if mrcore.FactoryErrStorageRowsNotAffected.Is(err) {
        return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, entityName)
    }

    return mrcore.FactoryErrServiceEntityNotRemoved.Caller(h.callerSkip).Wrap(err, entityName)
}

func (h *ServiceHelper) ReturnErrorIfItemNotFound(err error, entityName string) error {
    if err != nil {
        if mrcore.FactoryErrStorageNoRowFound.Is(err) {
            return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, entityName)
        }

        return mrcore.FactoryErrServiceTemporarilyUnavailable.Caller(h.callerSkip).Wrap(err, entityName)
    }

    return nil
}
