package mrjson

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// JsonDecoder - comment struct.
	JsonDecoder struct{}
)

// NewDecoder - создаёт объект JsonDecoder.
func NewDecoder() *JsonDecoder {
	return &JsonDecoder{}
}

// ParseToStruct - comment method.
func (p *JsonDecoder) ParseToStruct(_ context.Context, content io.Reader, structPointer any) error {
	dec := json.NewDecoder(content)
	dec.DisallowUnknownFields()

	if err := dec.Decode(structPointer); err != nil {
		return mrcore.ErrHttpRequestParseData.New(err.Error())
	}

	return nil
}
