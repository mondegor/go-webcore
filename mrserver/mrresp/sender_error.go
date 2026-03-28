package mrresp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlocale/model"
	"github.com/mondegor/go-sysmess/mrlog"
	mrtracectx "github.com/mondegor/go-sysmess/mrtrace/context"
	"github.com/mondegor/go-sysmess/util/xio"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// ErrorSender - формирует и отправляет клиенту ответ об ошибке.
	ErrorSender struct {
		encoder        mrserver.ResponseEncoder
		errorHandler   errors.Handler
		extractErrorID func(err error) string
		logger         mrlog.Logger
		parserLocale   parserLocale
		statusMapper   mrserver.ErrorStatusMapper
		debugInfo      func(err error) string
	}

	parserLocale interface {
		Localizer(r *http.Request) mrcore.Localizer
	}
)

// NewErrorSender - создаёт объект ErrorSender.
func NewErrorSender(
	encoder mrserver.ResponseEncoder,
	errorHandler errors.Handler,
	extractErrorID func(err error) string,
	logger mrlog.Logger,
	parserLocale parserLocale,
	statusMapper mrserver.ErrorStatusMapper,
	debugInfo func(err error) string,
) *ErrorSender {
	return &ErrorSender{
		encoder:        encoder,
		errorHandler:   errorHandler,
		extractErrorID: extractErrorID,
		logger:         logger,
		parserLocale:   parserLocale,
		statusMapper:   statusMapper,
		debugInfo:      debugInfo,
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

	if customError, ok := err.(errors.CustomError); ok { //nolint:errorlint
		if !customError.IsKindUser() {
			err = customError.Unwrap() // здесь возвращается ошибка с причиной невалидности

			goto InternalErrorSection
		}

		rs.errorHandler.Handle(ctx, customError.Unwrap()) // логируется пользовательская ошибка

		sendResponse(
			http.StatusBadRequest,
			rs.getErrorListResponse(r, customError),
		)

		return // OK, пользовательская ошибка обработана
	}

	if customErrors, ok := err.(errors.CustomListError); ok { //nolint:errorlint
		hasInvalidError := false

		for _, customError := range customErrors {
			if hasInvalidError || customError.IsKindUser() {
				rs.errorHandler.Handle(ctx, customError.Unwrap()) // логируются валидные пользовательские ошибки и невалидные начиная со второй

				continue
			}

			hasInvalidError = true
			err = customError.Unwrap() // здесь возвращается ошибка с причиной невалидности
		}

		if hasInvalidError {
			goto InternalErrorSection
		}

		sendResponse(
			http.StatusBadRequest,
			rs.getErrorListResponse(r, customErrors...),
		)

		return // OK, список пользовательских ошибок обработан
	}

InternalErrorSection:
	// в эту секцию поступают следующие виды ошибок:
	// 1. runtime ошибки;
	// 2. остальные ошибки у которых нет метода Kind() (требуется найти место их возникновения и правильно обработать);
	rs.errorHandler.Handle(ctx, err)

	status := rs.statusMapper.ErrorStatus(err)

	// TODO: подумать, нужно ли извлекать Code() из err для NewErrorAttribute

	if status == http.StatusBadRequest {
		sendResponse(
			status,
			[]ErrorAttribute{
				NewErrorAttribute(
					"",
					rs.parserLocale.Localizer(r).TranslateError(err),
					rs.debugInfo(err),
				),
			},
		)

		return
	}

	sendResponse(
		status,
		rs.getErrorDetailsResponse(r, err),
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

		rs.logger.Error(ctx, "marshal failed", "error", errors.ErrInternalHttpResponseParseData.Wrap(err))
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	xio.Write(ctx, rs.logger, w, bytes)
}

func (rs *ErrorSender) getErrorListResponse(r *http.Request, errorList ...errors.CustomError) ErrorListResponse {
	lz := rs.parserLocale.Localizer(r)
	attrs := make([]ErrorAttribute, len(errorList))

	for i, customError := range errorList {
		err := customError.Unwrap()
		attrs[i] = NewErrorAttribute(
			customError.CustomCode(),
			lz.TranslateError(err),
			rs.debugInfo(err),
		)
	}

	return attrs
}

func (rs *ErrorSender) getErrorDetailsResponse(r *http.Request, err error) ErrorDetailsResponse {
	errorMessage := model.ParseErrorMessage(rs.parserLocale.Localizer(r).TranslateError(err))
	errorTraceID := mrtracectx.RequestID(r.Context())

	if id := rs.extractErrorID(err); id != "" {
		errorTraceID = id
	}

	response := ErrorDetailsResponse{
		Title:        errorMessage.Reason,
		Details:      errorMessage.Details,
		Request:      r.Method + " " + r.URL.Path,
		Time:         time.Now().UTC().Format(time.RFC3339),
		ErrorTraceID: errorTraceID,
	}

	if debugInfo := rs.debugInfo(err); debugInfo != "" {
		if response.Details != "" {
			response.Details += "\n"
		}

		response.Details += "DebugInfo: " + debugInfo
	}

	return response
}
