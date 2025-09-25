package mrcore

type (
	// Localizer - comment interface.
	Localizer interface {
		Language() string
		Translate(message string, args ...any) string
		TranslateError(err error) string
		TranslateInDomain(domain, message string, args ...any) string
	}
)
