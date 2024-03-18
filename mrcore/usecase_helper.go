package mrcore

type (
	UsecaseHelper struct {
		callerSkipFrame int
	}
)

func NewUsecaseHelper() *UsecaseHelper {
	return &UsecaseHelper{
		callerSkipFrame: 2, // skip: wrapErrorFailed() + parent function
	}
}

func (h *UsecaseHelper) Caller(skipFrame int) *UsecaseHelper {
	if skipFrame == 0 {
		return h
	}

	c := *h
	c.callerSkipFrame += skipFrame

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
	factory := FactoryErrUseCaseTemporarilyUnavailable

	if FactoryErrStorageQueryFailed.Is(err) {
		factory = FactoryErrUseCaseOperationFailed
	}

	return factory.WithAttr(name, data).WithCaller(h.callerSkipFrame).Wrap(err)
}
