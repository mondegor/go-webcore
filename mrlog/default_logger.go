package mrlog

var def = New(InfoLevel).With().Str("mrlogger", "DEFAULT").Logger()

// Default - возвращает логгер по умолчанию (использовать только на крайней случай).
func Default() Logger {
	return def
}
