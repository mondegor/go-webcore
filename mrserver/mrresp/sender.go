package mrresp

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	Sender struct {
		encoder mrserver.ResponseEncoder
	}
)

// Make sure the Sender conforms with the mrserver.ResponseSender interface
var _ mrserver.ResponseSender = (*Sender)(nil)

func NewSender(encoder mrserver.ResponseEncoder) *Sender {
	return &Sender{
		encoder: encoder,
	}
}

func (rs *Sender) Send(w http.ResponseWriter, status int, structure any) error {
	bytes, err := rs.encoder.Marshal(structure)
	if err != nil {
		return mrcore.FactoryErrHTTPResponseParseData.Wrap(err)
	}

	rs.sendResponse(w, status, rs.encoder.ContentType(), bytes)

	return nil
}

func (rs *Sender) SendBytes(w http.ResponseWriter, status int, body []byte) error {
	rs.sendResponse(w, status, rs.encoder.ContentType(), body)

	return nil
}

func (rs *Sender) SendNoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (rs *Sender) sendResponse(w http.ResponseWriter, status int, contentType string, body []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)

	if len(body) < 1 {
		return
	}

	w.Write(body)
}
