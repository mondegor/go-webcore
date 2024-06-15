package mrdebug

import (
	"errors"
	"sync/atomic"

	"github.com/mondegor/go-webcore/mrcore"
)

var (
	once    atomic.Bool
	isDebug bool
)

// IsDebug  - comment func.
func IsDebug() bool {
	return isDebug
}

// EnableDebug  - comment func.
func EnableDebug() error {
	if once.Swap(true) {
		return mrcore.ErrInternal.Wrap(errors.New("debug flag is already enabled"))
	}

	isDebug = true

	return nil
}
