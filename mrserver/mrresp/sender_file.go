package mrresp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrmodel/media"
)

type (
	// FileSender - формирует и отправляет клиенту ответ с указанным файлом.
	// Наследует базовый Sender для отправки обычных ответов.
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

// SendFile - отправляет файл для отображения в браузере (inline).
// Content-Type устанавливается из FileInfo.ContentType файла.
// Content-Length устанавливается если известен размер файла.
func (rs *FileSender) SendFile(ctx context.Context, w http.ResponseWriter, file media.File) error {
	return rs.sendFile(ctx, w, file, false)
}

// SendAttachmentFile - отправляет файл как вложение для скачивания (attachment).
// Устанавливает заголовок Content-Disposition: attachment с именем файла.
// Устанавливает Cache-control: private для предотвращения кэширования прокси.
func (rs *FileSender) SendAttachmentFile(ctx context.Context, w http.ResponseWriter, file media.File) error {
	return rs.sendFile(ctx, w, file, true)
}

func (rs *FileSender) sendFile(ctx context.Context, w http.ResponseWriter, file media.File, isAttachment bool) error {
	w.Header().Set("Content-Type", file.ContentType)

	if file.Size > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))
	}

	if isAttachment {
		w.Header().Set("Cache-control", "private")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Original())) // TODO: escape
	}

	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, file.Body); err != nil {
		rs.logger.Error(ctx, "error file", "file", file.Path, "error", err)
	}

	return nil
}
