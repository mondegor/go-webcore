package mrcore

type (
	// Localizer - предоставляет методы для локализации сообщений и ошибок
	// с поддержкой перевода в разных доменах и форматирования аргументов.
	Localizer interface {
		// Language - возвращает текущий язык локализатора.
		Language() string

		// Translate - переводит сообщение на текущий язык с подстановкой аргументов.
		Translate(message string, args ...any) string

		// TranslateError - переводит ошибку на текущий язык локализатора.
		TranslateError(err error) string

		// TranslateInDomain - переводит сообщение в указанном домене с подстановкой аргументов.
		// Домен позволяет группировать переводы по контексту.
		TranslateInDomain(domain, message string, args ...any) string
	}
)
