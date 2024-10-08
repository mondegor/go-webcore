package mrlog

var def = New(InfoLevel).With().Str("mrlogger", "DEFAULT").Logger() //nolint:gochecknoglobals

// Default - возвращает логгер по умолчанию (использовать только на крайней случай).
func Default() Logger {
	return def
}
