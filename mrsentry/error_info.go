package mrsentry

type (
	errorInfo interface {
		Code() string
		Unwrap() error
	}
)

func unwrapErrorInfo(err errorInfo) errorInfo {
	unwrappedErr := err.Unwrap()

	if unwrappedErr == nil {
		return nil
	}

	err, ok := unwrappedErr.(errorInfo)
	if !ok {
		return nil
	}

	return err
}
