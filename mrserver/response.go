package mrserver

import (
	"context"
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/errors/helper"
	"github.com/mondegor/go-sysmess/mrmodel"
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
		SendFile(ctx context.Context, w http.ResponseWriter, file mrmodel.File) error
		SendAttachmentFile(ctx context.Context, w http.ResponseWriter, file mrmodel.File) error
	}

	// ErrorResponseSender - отправляет ответ со списком ошибок полученных в результате обработки запроса.
	ErrorResponseSender interface {
		SendError(w http.ResponseWriter, r *http.Request, err error)
	}

	// ErrorStatusMapper - возвращает http статус на основе указанной ошибки.
	ErrorStatusMapper interface {
		ErrorStatus(err error) int
	}
)

// NewHttpErrorStatusMapper - создаёт объект ErrorStatusMapper.
// Только для: 4XX, 5XX.
func NewHttpErrorStatusMapper(unexpectedStatus int, codeStatus ...any) (ErrorStatusMapper, error) {
	if unexpectedStatus <= 0 {
		unexpectedStatus = http.StatusInternalServerError
	}

	codeStatus = append(
		[]any{
			errors.ErrHttpClientUnauthorized.Code(), http.StatusUnauthorized,
			errors.ErrHttpAccessForbidden.Code(), http.StatusForbidden,
			errors.ErrAccessForbidden.Code(), http.StatusForbidden,
			errors.ErrHttpResourceNotFound.Code(), http.StatusNotFound,
			errors.ErrRecordNotFound.Code(), http.StatusNotFound,
			errors.ErrRecordVersionConflict.Code(), http.StatusConflict,
			errors.ErrHttpRequestParseData.Code(), http.StatusUnprocessableEntity,
			errors.ErrHttpTooManyRequests.Code(), http.StatusTooManyRequests,
			errors.ErrNotImplemented.Code(), http.StatusNotImplemented,
		},
		codeStatus...,
	)

	return helper.NewErrorStatusMapper(
		http.StatusBadRequest,
		http.StatusServiceUnavailable,
		http.StatusInternalServerError,
		unexpectedStatus,
		codeStatus,
	)
}
