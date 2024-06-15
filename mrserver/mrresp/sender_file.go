package mrresp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// FileSender - comment struct.
	FileSender struct {
		*Sender
	}
)

// Make sure the FileSender conforms with the mrserver.FileResponseSender interface.
var _ mrserver.FileResponseSender = (*FileSender)(nil)

// NewFileSender - создаёт объект FileSender.
func NewFileSender(base *Sender) *FileSender {
	return &FileSender{
		Sender: base,
	}
}

// SendFile - comment method.
func (rs *FileSender) SendFile(ctx context.Context, w http.ResponseWriter, file mrtype.File) error {
	return rs.sendFile(ctx, w, file, false)
}

// SendAttachmentFile - comment method.
func (rs *FileSender) SendAttachmentFile(ctx context.Context, w http.ResponseWriter, file mrtype.File) error {
	return rs.sendFile(ctx, w, file, true)
}

func (rs *FileSender) sendFile(ctx context.Context, w http.ResponseWriter, file mrtype.File, isAttachment bool) error {
	w.Header().Set("Content-Type", file.FileInfo.ContentType)

	if file.FileInfo.Size > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(file.FileInfo.Size, 10))
	}

	if isAttachment {
		w.Header().Set("Cache-control", "private")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileInfo.Original())) // TODO: escape
	}

	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, file.Body); err != nil {
		mrlog.Ctx(ctx).Error().Err(err).Msgf("error file %s", file.FileInfo.Path)
	}

	return nil
}
