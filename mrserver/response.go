package mrserver

import (
	"context"
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/errors/helper"
	"github.com/mondegor/go-sysmess/mrtype"
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

	// ErrorStatusMapper - возвращает http статус на основе указанной ошибки.
	ErrorStatusMapper interface {
		ErrorStatus(err error) int
	}
)

// NewHttpErrorStatusMapper - создаёт объект ErrorStatusMapper.
// Только для: 4XX, 5XX.
func NewHttpErrorStatusMapper(unexpectedStatus int, codeStatus ...any) ErrorStatusMapper {
	if unexpectedStatus == 0 {
		unexpectedStatus = http.StatusInternalServerError
	}

	codeStatus = append(
		[]any{
			errors.ErrHttpClientUnauthorized.Code(), http.StatusUnauthorized,
			errors.ErrHttpAccessForbidden.Code(), http.StatusForbidden,
			errors.ErrHttpResourceNotFound.Code(), http.StatusNotFound,
			errors.ErrUseCaseEntityNotFound.Code(), http.StatusNotFound,
			errors.ErrUseCaseEntityVersionConflict.Code(), http.StatusConflict,
			errors.ErrHttpRequestParseData.Code(), http.StatusUnprocessableEntity,
			errors.ErrHttpTooManyRequests.Code(), http.StatusTooManyRequests,
		},
		codeStatus...,
	)

	s, err := helper.NewErrorStatusMapper(
		http.StatusBadRequest,
		http.StatusServiceUnavailable,
		http.StatusInternalServerError,
		unexpectedStatus,
		codeStatus,
	)
	if err != nil {
		panic(err)
	}

	return s
}
