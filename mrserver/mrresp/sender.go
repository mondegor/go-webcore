package mrresp

import (
	"net/http"

	"github.com/mondegor/go-sysmess/errors"

	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// Sender - формирует и отправляет клиенту успешные HTTP-ответы.
	Sender struct {
		encoder mrserver.ResponseEncoder
	}
)

// NewSender - создаёт объект Sender.
func NewSender(encoder mrserver.ResponseEncoder) *Sender {
	return &Sender{
		encoder: encoder,
	}
}

// Send - отправляет клиенту ответ с сериализованными данными.
//
// Параметры:
//   - w - HTTP-ответ для записи;
//   - status - HTTP-код ответа (должен быть 2XX или 3XX);
//   - structure - данные для сериализации и отправки.
func (rs *Sender) Send(w http.ResponseWriter, status int, structure any) error {
	bytes, err := rs.encoder.Marshal(structure)
	if err != nil {
		return errors.ErrInternalHttpResponseParseData.Wrap(err)
	}

	return rs.sendResponse(w, status, rs.encoder.ContentType(), bytes)
}

// SendBytes - отправляет клиенту ответ в виде заранее подготовленных байтов.
func (rs *Sender) SendBytes(w http.ResponseWriter, status int, body []byte) error {
	return rs.sendResponse(w, status, rs.encoder.ContentType(), body)
}

// SendNoContent - отправляет ответ без содержимого со статусом 204 No Content.
func (rs *Sender) SendNoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (rs *Sender) sendResponse(w http.ResponseWriter, status int, contentType string, body []byte) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)

	if len(body) == 0 {
		return nil
	}

	_, err := w.Write(body)

	return err
}
