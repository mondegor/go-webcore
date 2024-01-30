package mrdebug

var (
	isDebug = false
)

func IsDebug() bool {
	return isDebug
}

// SetDebugFlag - WARNING: use only when starting the main process
func SetDebugFlag(value bool) {
	isDebug = value
}
