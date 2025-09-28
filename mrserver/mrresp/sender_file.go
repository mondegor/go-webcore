package mrresp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"
)

type (
	// FileSender - формирует и отправляет клиенту ответ с указанным файлом.
	FileSender struct {
		*Sender
		logger mrlog.Logger
	}
)

// NewFileSender - создаёт объект FileSender.
func NewFileSender(base *Sender, logger mrlog.Logger) *FileSender {
	return &FileSender{
		Sender: base,
		logger: logger,
	}
}

// SendFile - отправляет указанный файл, в случае неудачи возвращает ошибку.
func (rs *FileSender) SendFile(ctx context.Context, w http.ResponseWriter, file mrtype.File) error {
	return rs.sendFile(ctx, w, file, false)
}

// SendAttachmentFile - отправляет указанный файл в виде вложения для сохранения локально, в случае неудачи возвращает ошибку.
func (rs *FileSender) SendAttachmentFile(ctx context.Context, w http.ResponseWriter, file mrtype.File) error {
	return rs.sendFile(ctx, w, file, true)
}

func (rs *FileSender) sendFile(ctx context.Context, w http.ResponseWriter, file mrtype.File, isAttachment bool) error {
	w.Header().Set("Content-Type", file.FileInfo.ContentType)

	if file.FileInfo.Size > 0 {
		w.Header().Set("Content-Length", strconv.FormatUint(file.FileInfo.Size, 10))
	}

	if isAttachment {
		w.Header().Set("Cache-control", "private")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileInfo.Original())) // TODO: escape
	}

	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, file.Body); err != nil {
		rs.logger.Error(ctx, "error file", "file", file.FileInfo.Path, "error", err)
	}

	return nil
}
