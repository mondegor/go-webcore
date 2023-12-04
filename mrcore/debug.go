package mrcore

var (
	debug = false
)

func Debug() bool {
	return debug
}

// SetDebug - WARNING: use only when starting the main process
func SetDebug(value bool) {
	debug = value
}
