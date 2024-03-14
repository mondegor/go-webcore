package mrcore

type (
	UsecaseHelper struct {
		callerSkip int
	}
)

func NewUsecaseHelper() *UsecaseHelper {
	return &UsecaseHelper{
		callerSkip: 2, // skip: wrapErrorFailed and parent function
	}
}

func (h *UsecaseHelper) Caller(skip int) *UsecaseHelper {
	if skip == 0 {
		return h
	}

	c := *h
	c.callerSkip += skip

	return &c
}

func (h *UsecaseHelper) IsNotFoundError(err error) bool {
	return FactoryErrStorageNoRowFound.Is(err) ||
		FactoryErrStorageRowsNotAffected.Is(err)
}

func (h *UsecaseHelper) WrapErrorFailed(err error, source string) error {
	return h.wrapErrorFailed(err, "source", source)
}

func (h *UsecaseHelper) WrapErrorEntityFailed(err error, entityName string, entityData any) error {
	return h.wrapErrorFailed(err, entityName, entityData)
}

func (h *UsecaseHelper) WrapErrorEntityNotFoundOrFailed(err error, entityName string, entityData any) error {
	if h.IsNotFoundError(err) {
		return FactoryErrUseCaseEntityNotFound.Wrap(err)
	}

	return h.wrapErrorFailed(err, entityName, entityData)
}

func (h *UsecaseHelper) wrapErrorFailed(err error, name string, data any) error {
	if FactoryErrStorageQueryFailed.Is(err) {
		return FactoryErrUseCaseOperationFailed.WithAttr(name, data).Wrap(err)
	}

	return FactoryErrUseCaseTemporarilyUnavailable.WithAttr(name, data).Caller(h.callerSkip).Wrap(err)
}
