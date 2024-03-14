package mrjson

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
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
		mrlog.Ctx(ctx).Warn().Caller(1).Err(err).Send()
		return mrcore.FactoryErrHttpRequestParseData.Wrap(err)
	}

	return nil
}
