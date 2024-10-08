package mrinit

import "github.com/mondegor/go-webcore/mrcore"

// EnrichedInternalErrors - возвращает список ошибок наделённых дополнительными
// параметрами используемые в разных слоях приложения.
func EnrichedInternalErrors() []EnrichedError {
	return []EnrichedError{
		WrapProtoWithoutCaller(mrcore.ErrUnexpectedInternal),
		WrapProto(mrcore.ErrInternal),
		WrapProto(mrcore.ErrInternalNilPointer),
		WrapProto(mrcore.ErrInternalTypeAssertion),
		WrapProto(mrcore.ErrInternalInvalidType),
		WrapProto(mrcore.ErrInternalKeyNotFoundInSource),
		WrapProto(mrcore.ErrInternalFailedToOpen),
		WrapProto(mrcore.ErrInternalFailedToClose),
		WrapProto(mrcore.ErrInternalValueLenMax),
		WrapProto(mrcore.ErrInternalValueNotMatchRegexpPattern),
	}
}

// EnrichedStorageErrors - возвращает список ошибок наделённых дополнительными
// параметрами используемые при работе с хранилищами данных.
func EnrichedStorageErrors() []EnrichedError {
	return []EnrichedError{
		WrapProto(mrcore.ErrStorageConnectionIsAlreadyCreated),
		WrapProto(mrcore.ErrStorageConnectionIsNotOpened),
		WrapProto(mrcore.ErrStorageConnectionFailed),
		WrapProto(mrcore.ErrStorageQueryFailed),
		WrapProto(mrcore.ErrStorageFetchDataFailed),
		WrapProtoExtraDisabled(mrcore.ErrStorageNoRowFound),
		WrapProtoExtraDisabled(mrcore.ErrStorageRowsNotAffected),
	}
}

// EnrichedUseCaseErrors - возвращает список ошибок наделённых дополнительными
// параметрами используемые в бизнес-логике приложения.
func EnrichedUseCaseErrors() []EnrichedError {
	return []EnrichedError{
		WrapProto(mrcore.ErrUseCaseOperationFailed),
		WrapProto(mrcore.ErrUseCaseTemporarilyUnavailable),
		WrapProtoWithoutOnCreated(mrcore.ErrUseCaseIncorrectInputData),
		WrapProto(mrcore.ErrUseCaseEntityNotFound),
		WrapProto(mrcore.ErrUseCaseEntityNotAvailable),
		WrapProto(mrcore.ErrUseCaseEntityVersionInvalid),
		WrapProto(mrcore.ErrUseCaseSwitchStatusRejected),
		WrapProto(mrcore.ErrUseCaseInvalidFile),
	}
}

// EnrichedHttpErrors - возвращает список ошибок наделённых дополнительными
// параметрами используемые в http обработчиках.
func EnrichedHttpErrors() []EnrichedError {
	return []EnrichedError{
		WrapProto(mrcore.ErrHttpResponseParseData),
		WrapProto(mrcore.ErrHttpFileUpload),
		WrapProto(mrcore.ErrHttpMultipartFormFile),
		WrapProto(mrcore.ErrHttpClientUnauthorized),
		WrapProto(mrcore.ErrHttpAccessForbidden),
		WrapProto(mrcore.ErrHttpResourceNotFound),
		WrapProtoExtraDisabled(mrcore.ErrHttpRequestParseData),
		WrapProto(mrcore.ErrHttpRequestParseParam),
		WrapProto(mrcore.ErrHttpRequestParamEmpty),
		WrapProto(mrcore.ErrHttpRequestParamMax),
		WrapProto(mrcore.ErrHttpRequestParamLenMax),
	}
}
