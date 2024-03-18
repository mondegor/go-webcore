package mrjson

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	JsonDecoder struct {
	}
)

func NewDecoder() *JsonDecoder {
	return &JsonDecoder{}
}

func (p *JsonDecoder) ParseToStruct(ctx context.Context, content io.Reader, structPointer any) error {
	dec := json.NewDecoder(content)
	dec.DisallowUnknownFields()

	if err := dec.Decode(structPointer); err != nil {
		return mrcore.FactoryErrHttpRequestParseData.WithCaller(1).Wrap(err)
	}

	return nil
}
