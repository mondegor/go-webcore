package validate

import (
	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

type (
	// RequestContextParser - агрегирующий интерфейс парсеров контекста HTTP-запроса:
	// клиент, пользователь, локаль и часовой пояс.
	RequestContextParser interface {
		request.ParserClient
		request.ParserUser
		request.ParserLocale
		request.ParserTimeZone
	}

	// ContextParser - реализация RequestContextParser на основе парсеров go-webcore.
	ContextParser struct {
		*parser.Client
		*parser.User
		*parser.Locale
		*parser.TimeZone
	}
)

// NewContextParser - создаёт объект ContextParser.
func NewContextParser(
	clientParser *parser.Client,
	userParser *parser.User,
	localeParser *parser.Locale,
	timeZoneParser *parser.TimeZone,
) *ContextParser {
	return &ContextParser{
		Client:   clientParser,
		User:     userParser,
		Locale:   localeParser,
		TimeZone: timeZoneParser,
	}
}
