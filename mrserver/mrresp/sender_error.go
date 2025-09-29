package mrresp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlib/extio"
	"github.com/mondegor/go-sysmess/mrlocale/model"
	"github.com/mondegor/go-sysmess/mrlog"
	mrtracectx "github.com/mondegor/go-sysmess/mrtrace/context"

	core "github.com/mondegor/go-webcore/internal"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// ErrorSender - формирует и отправляет клиенту ответ об ошибке.
	ErrorSender struct {
		encoder      mrserver.ResponseEncoder
		errorHandler core.ErrorHandler
		logger       mrlog.Logger
		parserLocale parserLocale
		statusGetter mrserver.ErrorStatusGetter
		isDebug      bool
	}

	parserLocale interface {
		Localizer(r *http.Request) mrcore.Localizer
	}
)

// NewErrorSender - создаёт объект ErrorSender.
func NewErrorSender(
	encoder mrserver.ResponseEncoder,
	errorHandler core.ErrorHandler,
	logger mrlog.Logger,
	parserLocale parserLocale,
	statusGetter mrserver.ErrorStatusGetter,
	isDebug bool,
) *ErrorSender {
	return &ErrorSender{
		encoder:      encoder,
		errorHandler: errorHandler,
		logger:       logger,
		parserLocale: parserLocale,
		statusGetter: statusGetter,
		isDebug:      isDebug,
	}
}

// SendError - отправляет клиенту ответ об ошибке с одним из статусов: 4XX, 5XX и её деталями.
func (rs *ErrorSender) SendError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	sendResponse := func(status int, response any) {
		rs.sendStructResponse(
			ctx,
			w,
			status,
			rs.encoder.ContentTypeProblem(),
			response,
		)
	}

	if customError, ok := err.(*mrerr.CustomError); ok { //nolint:errorlint
		if !customError.IsValid() {
			err = customError.Err()

			goto InternalErrorSection
		}

		rs.errorHandler.Handle(ctx, customError.Err()) // логируется пользовательская ошибка

		sendResponse(
			http.StatusBadRequest,
			rs.getErrorListResponse(r, customError),
		)

		return // OK, пользовательская ошибка обработана
	}

	if customErrorList, ok := err.(mrerr.CustomErrors); ok { //nolint:errorlint
		hasInvalidError := false

		for _, customError := range customErrorList {
			if hasInvalidError || customError.IsValid() {
				rs.errorHandler.Handle(ctx, customError.Err()) // логируются все ошибки кроме первой невалидной

				continue
			}

			hasInvalidError = true
			err = customError.Err() // здесь возвращается внутренняя или системная ошибка
		}

		if hasInvalidError {
			goto InternalErrorSection
		}

		sendResponse(
			http.StatusBadRequest,
			rs.getErrorListResponse(r, customErrorList...),
		)

		return // OK, список пользовательских ошибок обработан
	}

InternalErrorSection:

	// в эту секцию поступают следующие виды ошибок:
	// 1. ошибки: InstantError kind=User/Internal/System;
	// 2. ProtoError kind=User/Internal/System (требуется найти место их создания и добавить для них вызов одного из методов New/Wrap);
	// 3. остальные ошибки: которые не были обёрнуты в InstantError (требуется найти место их создания и вложить их в InstantError);
	rs.errorHandler.HandleWith(
		ctx,
		err,
		func(analyzedKind mrerr.ErrorKind, err error) {
			status := rs.statusGetter.ErrorStatus(analyzedKind, err)

			if status == http.StatusBadRequest {
				sendResponse(
					status,
					[]ErrorAttribute{
						NewErrorAttribute(rs.parserLocale.Localizer(r), err, rs.isDebug),
					},
				)

				return
			}

			sendResponse(
				status,
				rs.getErrorDetailsResponse(r, err),
			)
		},
	)
}

func (rs *ErrorSender) sendStructResponse(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	contentType string,
	structResponse any,
) {
	bytes, err := json.Marshal(structResponse)
	if err != nil {
		status = http.StatusUnprocessableEntity
		bytes = nil

		rs.logger.Error(ctx, "marshal failed", "error", mr.ErrHttpResponseParseData.Wrap(err))
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	extio.Write(ctx, rs.logger, w, bytes)
}

func (rs *ErrorSender) getErrorListResponse(r *http.Request, errors ...*mrerr.CustomError) ErrorListResponse {
	lz := rs.parserLocale.Localizer(r)
	attrs := make([]ErrorAttribute, len(errors))

	for i, customError := range errors {
		attrs[i] = NewErrorAttribute(lz, customError, rs.isDebug)
	}

	return attrs
}

func (rs *ErrorSender) getErrorDetailsResponse(r *http.Request, err error) ErrorDetailsResponse {
	errorMessage := model.ParseErrorMessage(rs.parserLocale.Localizer(r).TranslateError(err))
	errorTraceID := mrtracectx.RequestID(r.Context())

	if e, ok := err.(interface{ ID() string }); ok && e.ID() != "" {
		errorTraceID = e.ID()
	}

	response := ErrorDetailsResponse{
		Title:        errorMessage.Reason,
		Details:      errorMessage.Details,
		Request:      r.Method + " " + r.URL.Path,
		Time:         time.Now().UTC().Format(time.RFC3339),
		ErrorTraceID: errorTraceID,
	}

	if rs.isDebug {
		if response.Details != "" {
			response.Details += "\n"
		}

		response.Details += "DebugInfo: " + err.Error()
	}

	return response
}
