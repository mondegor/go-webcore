package mrjson

import (
	"encoding/json"
)

const (
	// contentTypeJSON - стандартный MIME-тип для JSON-ответов.
	contentTypeJSON = "application/json; charset=utf-8"

	// contentTypeProblemJSON - MIME-тип для ответов с ошибками (RFC 9457 Problem Details).
	contentTypeProblemJSON = "application/problem+json; charset=utf-8"
)

type (
	// JsonEncoder - кодировщик Go-структур в JSON-формат.
	// Используется для формирования HTTP-ответов с JSON-содержимым.
	JsonEncoder struct{}
)

// NewEncoder - создаёт кодировщик JSON.
func NewEncoder() *JsonEncoder {
	return &JsonEncoder{}
}

// ContentType - возвращает MIME-тип для стандартных JSON-ответов.
func (p *JsonEncoder) ContentType() string {
	return contentTypeJSON
}

// ContentTypeProblem - возвращает MIME-тип для ответов с описанием ошибок (RFC 9457).
func (p *JsonEncoder) ContentTypeProblem() string {
	return contentTypeProblemJSON
}

// Marshal - сериализует Go-структуру в JSON-байты.
// Использует стандартный json.Marshal из пакета encoding/json.
// Параметр structure - данные для сериализации (структура, мапа, слайс и т.д.).
func (p *JsonEncoder) Marshal(structure any) ([]byte, error) {
	return json.Marshal(structure)
}
