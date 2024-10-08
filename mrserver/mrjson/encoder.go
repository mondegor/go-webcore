package mrjson

import (
	"encoding/json"
)

const (
	contentTypeJSON        = "application/json; charset=utf-8"
	contentTypeProblemJSON = "application/problem+json; charset=utf-8"
)

type (
	// JsonEncoder - comment struct.
	JsonEncoder struct{}
)

// NewEncoder - создаёт объект JsonEncoder.
func NewEncoder() *JsonEncoder {
	return &JsonEncoder{}
}

// ContentType - comment method.
func (p *JsonEncoder) ContentType() string {
	return contentTypeJSON
}

// ContentTypeProblem - comment method.
func (p *JsonEncoder) ContentTypeProblem() string {
	return contentTypeProblemJSON
}

// Marshal - comment method.
func (p *JsonEncoder) Marshal(structure any) ([]byte, error) {
	return json.Marshal(structure)
}
