package mrjson

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mondegor/go-core/errors"
)

type (
	// JsonDecoder - декодировщик JSON-данных в Go-структуры.
	// Используется для парсинга тела входящих HTTP-запросов.
	JsonDecoder struct{}
)

// NewDecoder - создаёт декодировщик JSON.
func NewDecoder() *JsonDecoder {
	return &JsonDecoder{}
}

// ParseToStruct - декодирует JSON из reader в Go-структуру.
//
// Особенности:
//   - Использует DisallowUnknownFields() - запросы с неизвестными полями будут отклонены;
//   - Ошибки декодирования оборачиваются в ErrHttpRequestParseData для отправки клиенту;
//   - structPointer - должен быть указателем на структуру.
func (p *JsonDecoder) ParseToStruct(_ context.Context, content io.Reader, structPointer any) error {
	dec := json.NewDecoder(content)
	dec.DisallowUnknownFields()

	if err := dec.Decode(structPointer); err != nil {
		return errors.ErrHttpRequestParseData.New(err) // данная ошибка передаётся пользователю в виде сообщения
	}

	return nil
}
