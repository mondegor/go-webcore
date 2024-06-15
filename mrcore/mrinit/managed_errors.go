package mrinit

import "github.com/mondegor/go-webcore/mrcore"

// ManagedInternalErrors - comment func.
func ManagedInternalErrors() []ManagedError {
	return []ManagedError{
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

// ManagedStorageErrors - comment func.
func ManagedStorageErrors() []ManagedError {
	return []ManagedError{
		WrapProto(mrcore.ErrStorageConnectionIsAlreadyCreated),
		WrapProto(mrcore.ErrStorageConnectionIsNotOpened),
		WrapProto(mrcore.ErrStorageConnectionFailed),
		WrapProto(mrcore.ErrStorageQueryFailed),
		WrapProto(mrcore.ErrStorageFetchDataFailed),
		WrapProtoExtraDisabled(mrcore.ErrStorageNoRowFound),
		WrapProtoExtraDisabled(mrcore.ErrStorageRowsNotAffected),
	}
}

// ManagedUseCaseErrors - comment func.
func ManagedUseCaseErrors() []ManagedError {
	return []ManagedError{
		WrapProto(mrcore.ErrUseCaseOperationFailed),
		WrapProto(mrcore.ErrUseCaseTemporarilyUnavailable),
		WrapProtoExtraDisabled(mrcore.ErrUseCaseIncorrectInputData),
		WrapProto(mrcore.ErrUseCaseEntityNotFound),
		WrapProto(mrcore.ErrUseCaseEntityNotAvailable),
		WrapProto(mrcore.ErrUseCaseEntityVersionInvalid),
		WrapProto(mrcore.ErrUseCaseSwitchStatusRejected),
		WrapProto(mrcore.ErrUseCaseInvalidFile),
	}
}

// ManagedHttpErrors - comment func.
func ManagedHttpErrors() []ManagedError {
	return []ManagedError{
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
