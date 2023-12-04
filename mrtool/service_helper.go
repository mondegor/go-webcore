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
		callerSkip: 1,
	}
}

func (h *ServiceHelper) Caller(skip int) *ServiceHelper {
	return &ServiceHelper{
		callerSkip: h.callerSkip + skip,
	}
}

func (h *ServiceHelper) WrapErrorEntityFetch(err error, entityName string, entityData any) error {
	if mrcore.FactoryErrStorageNoRowFound.Is(err) {
		return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, mrerr.Arg{entityName: entityData})
	}

	return mrcore.FactoryErrServiceTemporarilyUnavailable.Caller(h.callerSkip).Wrap(err, mrerr.Arg{entityName: entityData})
}

func (h *ServiceHelper) WrapErrorEntityInsert(err error, entityName string, entityData any) error {
	if mrcore.FactoryErrStorageQueryFailed.Is(err) {
		return mrcore.FactoryErrServiceEntityNotCreated.Wrap(err, mrerr.Arg{entityName: entityData})
	}

	return mrcore.FactoryErrServiceTemporarilyUnavailable.Caller(h.callerSkip).Wrap(err, mrerr.Arg{entityName: entityData})
}

func (h *ServiceHelper) WrapErrorEntityUpdate(err error, entityName string, entityData any) error {
	if mrcore.FactoryErrStorageRowsNotAffected.Is(err) ||
		mrcore.FactoryErrStorageNoRowFound.Is(err) {
		return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, mrerr.Arg{entityName: entityData})
	}

	if mrcore.FactoryErrStorageQueryFailed.Is(err) {
		return mrcore.FactoryErrServiceEntityNotStored.Wrap(err, mrerr.Arg{entityName: entityData})
	}

	return mrcore.FactoryErrServiceTemporarilyUnavailable.Caller(h.callerSkip).Wrap(err, mrerr.Arg{entityName: entityData})
}

func (h *ServiceHelper) WrapErrorEntityDelete(err error, entityName string, entityData any) error {
	if mrcore.FactoryErrStorageRowsNotAffected.Is(err) ||
		mrcore.FactoryErrStorageNoRowFound.Is(err) {
		return mrcore.FactoryErrServiceEntityNotFound.Wrap(err, mrerr.Arg{entityName: entityData})
	}

	if mrcore.FactoryErrStorageQueryFailed.Is(err) {
		return mrcore.FactoryErrServiceEntityNotRemoved.Wrap(err, mrerr.Arg{entityName: entityData})
	}

	return mrcore.FactoryErrServiceTemporarilyUnavailable.Caller(h.callerSkip).Wrap(err, mrerr.Arg{entityName: entityData})
}
