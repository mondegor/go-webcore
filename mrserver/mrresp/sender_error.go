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
		debugFunc      func(value any) string
	}

	// parserLocale - внутренний интерфейс для получения локализатора из запроса.
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
	debugFunc func(value any) string,
) *ErrorSender {
	return &ErrorSender{
		encoder:        encoder,
		errorHandler:   errorHandler,
		extractErrorID: extractErrorID,
		logger:         logger,
		parserLocale:   parserLocale,
		statusMapper:   statusMapper,
		debugFunc:      debugFunc,
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
		err = customError.Unwrap() // здесь возвращается ошибка с причиной невалидности

		if !customError.IsKindUser() {
			goto GeneralErrorSection
		}

		rs.errorHandler.Handle(ctx, err) // логируется пользовательская ошибка

		sendResponse(
			http.StatusBadRequest,
			NewError400Response(
				r,
				ErrorAttribute{
					Code:      customError.CustomCode(),
					Detail:    rs.parserLocale.Localizer(r).TranslateError(err),
					DebugInfo: rs.debugFunc(err),
				},
			),
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
			goto GeneralErrorSection
		}

		sendResponse(
			http.StatusBadRequest,
			NewError400Response(
				r,
				rs.customErrorsToErrorAttrs(r, customErrors)...,
			),
		)

		return // OK, список пользовательских ошибок обработан
	}

GeneralErrorSection:
	// в эту секцию поступают следующие виды ошибок:
	// 1. runtime ошибки;
	// 2. остальные ошибки у которых нет метода Kind() (требуется найти место их возникновения и правильно обработать);
	rs.errorHandler.Handle(ctx, err)

	status := rs.statusMapper.ErrorStatus(err)

	if status == http.StatusBadRequest {
		errorCode := ErrorAttributeCodeByDefault

		if e, ok := err.(interface{ Code() string }); ok && e.Code() != "" {
			errorCode = e.Code()
		}

		sendResponse(
			status,
			NewError400Response(
				r,
				ErrorAttribute{
					Code:      errorCode,
					Detail:    rs.parserLocale.Localizer(r).TranslateError(err),
					DebugInfo: rs.debugFunc(err),
				},
			),
		)

		return
	}

	sendResponse(
		status,
		rs.getErrorDetailsResponse(r, status, err),
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

func (rs *ErrorSender) customErrorsToErrorAttrs(r *http.Request, errorList []errors.CustomError) []ErrorAttribute {
	lz := rs.parserLocale.Localizer(r)
	attrs := make([]ErrorAttribute, len(errorList))

	for i, customError := range errorList {
		err := customError.Unwrap()
		attrs[i] = ErrorAttribute{
			Code:      customError.CustomCode(),
			Detail:    lz.TranslateError(err),
			DebugInfo: rs.debugFunc(err),
		}
	}

	return attrs
}

func (rs *ErrorSender) getErrorDetailsResponse(r *http.Request, status int, err error) ErrorDetailsResponse {
	errorMessage := model.ParseErrorMessage(rs.parserLocale.Localizer(r).TranslateError(err))
	errorTraceID := mrtracectx.RequestID(r.Context())

	if id := rs.extractErrorID(err); id != "" {
		errorTraceID = id
	}

	response := ErrorDetailsResponse{
		Type:         errorMessage.ProblemURL,
		Title:        errorMessage.Reason,
		Status:       status,
		Detail:       errorMessage.Details,
		Instance:     r.Method + " " + r.URL.Path,
		Time:         time.Now().UTC().Format(time.RFC3339),
		ErrorTraceID: errorTraceID,
		DebugInfo:    rs.debugFunc(err),
	}

	return response
}
