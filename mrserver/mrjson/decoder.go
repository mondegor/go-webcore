package mrjson

import (
	"encoding/json"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

type (
	JsonDecoder struct {
	}
)

func NewDecoder() *JsonDecoder {
	return &JsonDecoder{}
}

func (p *JsonDecoder) ParseToStruct(r *http.Request, structPointer any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(structPointer); err != nil {
		mrctx.Logger(r.Context()).Caller(1).Warn(err)
		return mrcore.FactoryErrHttpRequestParseData.Wrap(err)
	}

	return nil
}