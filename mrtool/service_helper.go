package mrtool

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
)

type (
	ServiceHelper struct {
		callerSkip int
	}
)

func NewServiceHelper() *ServiceHelper {
	return &ServiceHelper{
		callerSkip: 2, // skip: wrapErrorFailed and parent function
	}
}

func (h *ServiceHelper) Caller(skip int) *ServiceHelper {
	if skip == 0 {
		return h
	}

	c := *h
	c.callerSkip += skip

	return &c
}

func (h *ServiceHelper) IsNotFound(err error) bool {
	return mrcore.FactoryErrStorageNoRowFound.Is(err) ||
		mrcore.FactoryErrStorageRowsNotAffected.Is(err)
}

func (h *ServiceHelper) WrapErrorFailed(err error, source string) error {
	return h.wrapErrorFailed(err, "source", source)
}

func (h *ServiceHelper) WrapErrorEntityFailed(err error, entityName string, entityData any) error {
	return h.wrapErrorFailed(err, entityName, entityData)
}

func (h *ServiceHelper) WrapErrorEntity(
	reasonWhyNotFound *mrerr.AppErrorFactory,
	err error,
	entityName string,
	entityData any,
) error {
	if h.IsNotFound(err) {
		return reasonWhyNotFound.Wrap(err)
	}

	return h.wrapErrorFailed(err, entityName, entityData)
}

func (h *ServiceHelper) WrapErrorEntityNotFoundOrFailed(err error, entityName string, entityData any) error {
	return h.WrapErrorEntity(
		mrcore.FactoryErrServiceEntityNotFound,
		err,
		entityName,
		entityData,
	)
}

func (h *ServiceHelper) wrapErrorFailed(err error, name string, data any) error {
	err = mrcore.FactoryErrWithData.Wrap(err, name, data)

	if mrcore.FactoryErrStorageQueryFailed.Is(err) {
		return mrcore.FactoryErrServiceOperationFailed.Wrap(err)
	}

	return mrcore.FactoryErrServiceTemporarilyUnavailable.Caller(h.callerSkip).Wrap(err)
}
