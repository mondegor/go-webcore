package mrjson

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// JsonDecoder - comment struct.
	JsonDecoder struct{}
)

// Make sure the Image conforms with the mrserver.RequestDecoder interface.
var _ mrserver.RequestDecoder = (*JsonDecoder)(nil)

// NewDecoder - создаёт объект JsonDecoder.
func NewDecoder() *JsonDecoder {
	return &JsonDecoder{}
}

// ParseToStruct - comment method.
func (p *JsonDecoder) ParseToStruct(_ context.Context, content io.Reader, structPointer any) error {
	dec := json.NewDecoder(content)
	dec.DisallowUnknownFields()

	if err := dec.Decode(structPointer); err != nil {
		return mrcore.ErrHttpRequestParseData.Wrap(err)
	}

	return nil
}
