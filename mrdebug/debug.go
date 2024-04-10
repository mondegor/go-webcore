package mrdebug

var (
	isDebug = false
)

func IsDebug() bool {
	return isDebug
}

// SetDebugFlag - WARNING: use only by the main process when it is starting
func SetDebugFlag(value bool) {
	isDebug = value
}
