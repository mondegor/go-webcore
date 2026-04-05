package mrjson

import (
	"encoding/json"
)

const (
	contentTypeJSON        = "application/json; charset=utf-8"
	contentTypeProblemJSON = "application/problem+json; charset=utf-8"
)

type (
	// JsonEncoder - кодирует Go-структуры в JSON-формат.
	JsonEncoder struct{}
)

// NewEncoder - создаёт объект JsonEncoder.
func NewEncoder() *JsonEncoder {
	return &JsonEncoder{}
}

// ContentType - возвращает MIME тип для JSON контента.
func (p *JsonEncoder) ContentType() string {
	return contentTypeJSON
}

// ContentTypeProblem - возвращает MIME тип для JSON контента с проблемой (application/problem+json).
func (p *JsonEncoder) ContentTypeProblem() string {
	return contentTypeProblemJSON
}

// Marshal - сериализует Go структуру в JSON байты.
func (p *JsonEncoder) Marshal(structure any) ([]byte, error) {
	return json.Marshal(structure)
}
