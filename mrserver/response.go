package mrserver

import (
	"context"
	"net/http"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// ResponseEncoder - формирует ответ сервера из go структуры в необходимом формате.
	ResponseEncoder interface {
		ContentType() string
		ContentTypeProblem() string
		Marshal(structure any) ([]byte, error)
	}

	// ResponseSender - отправляет ответ с данными сформированные сервером.
	ResponseSender interface {
		Send(w http.ResponseWriter, status int, structure any) error
		SendBytes(w http.ResponseWriter, status int, body []byte) error
		SendNoContent(w http.ResponseWriter) error
	}

	// FileResponseSender - отправляет ответ с данными в виде файла сформированные сервером.
	FileResponseSender interface {
		ResponseSender
		SendFile(ctx context.Context, w http.ResponseWriter, file mrtype.File) error
		SendAttachmentFile(ctx context.Context, w http.ResponseWriter, file mrtype.File) error
	}

	// ErrorResponseSender - отправляет ответ со списком ошибок полученных в результате обработки запроса.
	ErrorResponseSender interface {
		SendError(w http.ResponseWriter, r *http.Request, err error)
	}

	// ErrorStatusGetter - возвращает http статус на основе указанной ошибки.
	ErrorStatusGetter interface {
		ErrorStatus(err error) int
	}
)
