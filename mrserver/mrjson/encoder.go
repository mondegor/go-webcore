package mrjson

import "encoding/json"

const (
	contentTypeJSON        = "application/json; charset=utf-8"
	contentTypeProblemJSON = "application/problem+json; charset=utf-8"
)

type (
	JSONEncoder struct{}
)

func NewEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

func (p *JSONEncoder) ContentType() string {
	return contentTypeJSON
}

func (p *JSONEncoder) ContentTypeProblem() string {
	return contentTypeProblemJSON
}

func (p *JSONEncoder) Marshal(structure any) ([]byte, error) {
	return json.Marshal(structure)
}
