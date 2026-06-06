package mrserver

import (
	"context"
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/errors/helper"
	"github.com/mondegor/go-sysmess/mrmodel/media"
)

type (
	// ResponseEncoder - интерфейс кодировщика ответов сервера.
	// Преобразует Go-структуры в байтовый формат ответа (JSON, XML и т.д.).
	ResponseEncoder interface {
		ContentType() string
		ContentTypeProblem() string
		Marshal(structure any) ([]byte, error)
	}

	// ResponseSender - интерфейс отправки успешных ответов клиенту.
	ResponseSender interface {
		Send(w http.ResponseWriter, status int, structure any) error
		SendBytes(w http.ResponseWriter, status int, body []byte) error
		SendNoContent(w http.ResponseWriter) error
	}

	// FileResponseSender - интерфейс отправки файлов клиенту.
	FileResponseSender interface {
		ResponseSender

		SendFile(ctx context.Context, w http.ResponseWriter, file media.File) error
		SendAttachmentFile(ctx context.Context, w http.ResponseWriter, file media.File) error
	}

	// ErrorResponseSender - интерфейс отправки клиенту ответов с информацией об ошибках.
	ErrorResponseSender interface {
		SendError(w http.ResponseWriter, r *http.Request, err error)
	}

	// ErrorStatusMapper - интерфейс маппера ошибок в HTTP-статусы.
	ErrorStatusMapper interface {
		ErrorStatus(err error) int
	}
)

// NewHttpErrorStatusMapper - создаёт маппер ошибок в HTTP-статусы.
//
// Преднастроенные соответствия (могут быть переопределены через codeStatus):
//   - ErrHttpClientUnauthorized -> 401 Unauthorized
//   - ErrHttpAccessForbidden -> 403 Forbidden
//   - ErrAccessForbidden -> 403 Forbidden
//   - ErrHttpResourceNotFound -> 404 Not Found
//   - ErrRecordNotFound -> 404 Not Found
//   - ErrRecordVersionConflict -> 409 Conflict
//   - ErrHttpRequestParseData -> 422 Unprocessable Entity
//   - ErrHttpTooManyRequests -> 429 Too Many Requests
//   - ErrNotImplemented -> 501 Not Implemented
//
// Параметры:
//   - unexpectedStatus - HTTP-статус для непредвиденных ошибок (по умолчанию 500);
//   - codeStatus - пары "код ошибки -> HTTP-статус" для переопределения или добавления маппинга;
//
// Возвращает только статусы 4xx и 5xx.
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
