package mrresp

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// Sender - формирует и отправляет клиенту успешный ответ.
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

// Send - отправляет клиенту ответ с данными в виде структуры с одним из статусов: 2XX, 3XX.
func (rs *Sender) Send(w http.ResponseWriter, status int, structure any) error {
	bytes, err := rs.encoder.Marshal(structure)
	if err != nil {
		return mr.ErrHttpResponseParseData.Wrap(err)
	}

	return rs.sendResponse(w, status, rs.encoder.ContentType(), bytes)
}

// SendBytes - отправляет клиенту ответ у указанным массивом байт с одним из статусов: 2XX, 3XX.
func (rs *Sender) SendBytes(w http.ResponseWriter, status int, body []byte) error {
	return rs.sendResponse(w, status, rs.encoder.ContentType(), body)
}

// SendNoContent - отправляет клиенту ответ без данных со статусом 204.
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
