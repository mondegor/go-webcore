package mrjson

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mondegor/go-sysmess/errors"
)

type (
	// JsonDecoder - декодирует JSON-данные в Go-структуры.
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
		return errors.ErrHttpRequestParseData.New(err) // данная ошибка передаётся пользователю в виде сообщения
	}

	return nil
}
