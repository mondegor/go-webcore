package mrinit

import "github.com/mondegor/go-sysmess/mrerr"

type (
	// ManagedError - comment struct.
	ManagedError struct {
		Err           *mrerr.ProtoAppError
		WithCaller    bool
		WithOnCreated bool
	}
)

// WrapProto - comment func.
func WrapProto(proto *mrerr.ProtoAppError) ManagedError {
	if proto.Kind() == mrerr.ErrorKindInternal || proto.Kind() == mrerr.ErrorKindSystem {
		return ManagedError{
			Err:           proto,
			WithCaller:    true,
			WithOnCreated: true,
		}
	}

	return ManagedError{
		Err:           proto,
		WithCaller:    false,
		WithOnCreated: false,
	}
}

// WrapProtoList - comment func.
func WrapProtoList(protos []*mrerr.ProtoAppError) []ManagedError {
	errors := make([]ManagedError, len(protos))

	for i := range protos {
		errors[i] = WrapProto(protos[i])
	}

	return errors
}

// WrapProtoExtra - comment func.
func WrapProtoExtra(proto *mrerr.ProtoAppError, withCaller, withOnCreated bool) ManagedError {
	return ManagedError{
		Err:           proto,
		WithCaller:    withCaller,
		WithOnCreated: withOnCreated,
	}
}

// WrapProtoExtraEnabled - comment func.
func WrapProtoExtraEnabled(proto *mrerr.ProtoAppError) ManagedError {
	return ManagedError{
		Err:           proto,
		WithCaller:    true,
		WithOnCreated: true,
	}
}

// WrapProtoExtraDisabled - comment func.
func WrapProtoExtraDisabled(proto *mrerr.ProtoAppError) ManagedError {
	return ManagedError{
		Err:           proto,
		WithCaller:    false,
		WithOnCreated: false,
	}
}
