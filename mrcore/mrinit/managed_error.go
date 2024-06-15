package mrinit

import "github.com/mondegor/go-sysmess/mrerr"

type (
	// ManagedError - comment struct.
	ManagedError struct {
		Err             *mrerr.ProtoAppError
		WithIDGenerator bool
		WithCaller      bool
	}
)

// WrapProto - comment func.
func WrapProto(proto *mrerr.ProtoAppError) ManagedError {
	if proto.Kind() == mrerr.ErrorKindInternal || proto.Kind() == mrerr.ErrorKindSystem {
		return ManagedError{
			Err:             proto,
			WithIDGenerator: true,
			WithCaller:      true,
		}
	}

	return ManagedError{
		Err: proto,
	}
}

// WrapProtoList - comment func.
func WrapProtoList(protos []*mrerr.ProtoAppError) []ManagedError {
	errors := make([]ManagedError, 0, len(protos))

	for i := range protos {
		errors = append(errors, WrapProto(protos[i]))
	}

	return errors
}

// WrapProtoExtra - comment func.
func WrapProtoExtra(proto *mrerr.ProtoAppError, withIDGenerator, withCaller bool) ManagedError {
	return ManagedError{
		Err:             proto,
		WithIDGenerator: withIDGenerator,
		WithCaller:      withCaller,
	}
}

// WrapProtoExtraEnabled - comment func.
func WrapProtoExtraEnabled(proto *mrerr.ProtoAppError) ManagedError {
	return ManagedError{
		Err:             proto,
		WithIDGenerator: true,
		WithCaller:      true,
	}
}

// WrapProtoExtraDisabled - comment func.
func WrapProtoExtraDisabled(proto *mrerr.ProtoAppError) ManagedError {
	return ManagedError{
		Err:             proto,
		WithIDGenerator: false,
		WithCaller:      false,
	}
}
