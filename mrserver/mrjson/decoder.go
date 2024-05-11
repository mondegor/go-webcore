package mrjson

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	JSONDecoder struct{}
)

func NewDecoder() *JSONDecoder {
	return &JSONDecoder{}
}

func (p *JSONDecoder) ParseToStruct(ctx context.Context, content io.Reader, structPointer any) error {
	dec := json.NewDecoder(content)
	dec.DisallowUnknownFields()

	if err := dec.Decode(structPointer); err != nil {
		const skipFrame = 1
		return mrcore.FactoryErrHTTPRequestParseData.WithCallerSkipFrame(skipFrame).Wrap(err)
	}

	return nil
}
