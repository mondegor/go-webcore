package mrresponse

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileSender struct {
		*Sender
	}
)

// Make sure the FileSender conforms with the mrserver.FileResponseSender interface
var _ mrserver.FileResponseSender = (*FileSender)(nil)

func NewFileSender(base *Sender) *FileSender {
	return &FileSender{
		Sender: base,
	}
}

func (rs *FileSender) SendFile(w http.ResponseWriter, info mrtype.FileInfo, attachmentName string, file io.Reader) error {
	w.Header().Set("Content-Type", info.ContentType)

	if info.Size > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(info.Size, 10))
	}

	if attachmentName != "" {
		w.Header().Set("Cache-control", "private")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachmentName)) // :TODO: escape
	}

	w.WriteHeader(http.StatusOK)

	_, err := io.Copy(w, file)

	return err
}
