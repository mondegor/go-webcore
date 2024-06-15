package mrlog

import (
	"errors"
	"sync/atomic"

	"github.com/mondegor/go-webcore/mrcore"
)

var (
	once atomic.Bool
	def  Logger
)

// Default - возвращает ранее установленный логгер по умолчанию.
// Если ранее не был вызван метод mrlog.SetDefault(), то будет panic.
func Default() Logger {
	if def == nil {
		panic("no default logger, use: mrlog.SetDefault()")
	}

	return def
}

// SetDefault - устанавливает логгер по умолчанию, одноразовая операция.
func SetDefault(logger Logger) error {
	if def != nil {
		return mrcore.ErrInternal.Wrap(errors.New("default logger is already exists"))
	}

	if logger == nil {
		return mrcore.ErrInternalNilPointer.Wrap(errors.New("logger variable is a nil pointer, expected mrlog.Logger"))
	}

	if !once.Swap(true) {
		def = logger.With().Str("mrlogger", "DEFAULT").Logger()
	}

	return nil
}
