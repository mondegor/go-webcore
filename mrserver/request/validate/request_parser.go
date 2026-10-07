package validate

import (
	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

type (
	// RequestParser - агрегирующий интерфейс базовых парсеров HTTP-запроса: значения параметров и валидация тела.
	RequestParser interface {
		request.ParserBool
		request.ParserInt64
		request.ParserUint64
		request.ParserFloat64
		request.ParserString
		request.ParserDateTime
		request.ParserUUID
		request.ParserValidate
	}

	// Parser - реализация RequestParser на основе парсеров go-webcore.
	Parser struct {
		*parser.Bool
		*parser.Int64
		*parser.Uint64
		*parser.Float64
		*parser.String
		*parser.DateTime
		*parser.UUID
		*parser.Validator
	}
)

// NewParser - создаёт объект Parser.
func NewParser(
	boolParser *parser.Bool,
	int64Parser *parser.Int64,
	uint64Parser *parser.Uint64,
	float64Parser *parser.Float64,
	stringParser *parser.String,
	dateTimeParser *parser.DateTime,
	uuidParser *parser.UUID,
	validatorParser *parser.Validator,
) *Parser {
	return &Parser{
		Bool:      boolParser,
		Int64:     int64Parser,
		Uint64:    uint64Parser,
		Float64:   float64Parser,
		String:    stringParser,
		DateTime:  dateTimeParser,
		UUID:      uuidParser,
		Validator: validatorParser,
	}
}
