package mrcore

type (
	// Localizer - предоставляет методы для локализации сообщений и ошибок
	// с поддержкой перевода в разных доменах и форматирования аргументов.
	Localizer interface {
		Language() string
		Translate(message string, args ...any) string
		TranslateError(err error) string
		TranslateInDomain(domain, message string, args ...any) string
	}
)
