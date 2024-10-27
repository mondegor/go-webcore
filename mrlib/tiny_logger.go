package mrlib

type (
	logger interface {
		Printf(format string, args ...any)
	}
)
