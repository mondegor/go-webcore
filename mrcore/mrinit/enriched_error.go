package mrinit

import "github.com/mondegor/go-sysmess/mrerr"

type (
	// EnrichedError - обёртка ошибки ProtoAppError, наделяющая её дополнительными параметрами.
	EnrichedError struct {
		Err           *mrerr.ProtoAppError
		WithCaller    bool
		WithOnCreated bool
	}
)

// WrapProto - возвращает обёртку ошибки, которая:
// - для системных ошибок формирует стек вызовов и генерирует событие создания ошибки;
// - у пользовательских ошибок отключает всех эти опции.
func WrapProto(proto *mrerr.ProtoAppError) EnrichedError {
	if proto.Kind() == mrerr.ErrorKindInternal || proto.Kind() == mrerr.ErrorKindSystem {
		return EnrichedError{
			Err:           proto,
			WithCaller:    true,
			WithOnCreated: true,
		}
	}

	return EnrichedError{
		Err:           proto,
		WithCaller:    false,
		WithOnCreated: false,
	}
}

// WrapProtoList - возвращает массив обёрток, которые формируются при помощи WrapProto().
func WrapProtoList(protos []*mrerr.ProtoAppError) []EnrichedError {
	errors := make([]EnrichedError, len(protos))

	for i := range protos {
		errors[i] = WrapProto(protos[i])
	}

	return errors
}

// WrapProtoExtraEnabled - возвращает обёртку ошибки,
// которая всегда формирует стек вызовов и генерирует событие создания ошибки.
func WrapProtoExtraEnabled(proto *mrerr.ProtoAppError) EnrichedError {
	return EnrichedError{
		Err:           proto,
		WithCaller:    true,
		WithOnCreated: true,
	}
}

// WrapProtoWithoutCaller - возвращает обёртку ошибки,
// которая всегда генерирует событие создания ошибки и отключает формирование стека вызовов.
func WrapProtoWithoutCaller(proto *mrerr.ProtoAppError) EnrichedError {
	return EnrichedError{
		Err:           proto,
		WithCaller:    false,
		WithOnCreated: true,
	}
}

// WrapProtoWithoutOnCreated - возвращает обёртку ошибки,
// которая всегда формирует стек вызовов и отключает генерацию события создания ошибки.
func WrapProtoWithoutOnCreated(proto *mrerr.ProtoAppError) EnrichedError {
	return EnrichedError{
		Err:           proto,
		WithCaller:    true,
		WithOnCreated: false,
	}
}

// WrapProtoExtraDisabled - возвращает обёртку ошибки, которая отключает все опции.
func WrapProtoExtraDisabled(proto *mrerr.ProtoAppError) EnrichedError {
	return EnrichedError{
		Err:           proto,
		WithCaller:    false,
		WithOnCreated: false,
	}
}
