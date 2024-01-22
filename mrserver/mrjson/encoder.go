package mrjson

import "encoding/json"

const (
	contentTypeJson        = "application/json; charset=utf-8"
	contentTypeProblemJson = "application/problem+json; charset=utf-8"
)

type (
	JsonEncoder struct {
	}
)

func NewEncoder() *JsonEncoder {
	return &JsonEncoder{}
}

func (p *JsonEncoder) ContentType() string {
	return contentTypeJson
}

func (p *JsonEncoder) ContentTypeProblem() string {
	return contentTypeProblemJson
}

func (p *JsonEncoder) Marshal(structure any) ([]byte, error) {
	return json.Marshal(structure)
}
