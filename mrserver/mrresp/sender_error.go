package mrresp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlang"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// ErrorSender - формирует и отправляет клиенту ответ об ошибке.
	ErrorSender struct {
		encoder      mrserver.ResponseEncoder
		errorHandler mrcore.ErrorHandler
		statusGetter mrserver.ErrorStatusGetter
		isDebug      bool
	}
)

// NewErrorSender - создаёт объект ErrorSender.
func NewErrorSender(
	encoder mrserver.ResponseEncoder,
	errorHandler mrcore.ErrorHandler,
	statusGetter mrserver.ErrorStatusGetter,
	isDebug bool,
) *ErrorSender {
	return &ErrorSender{
		encoder:      encoder,
		errorHandler: errorHandler,
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

		sendResponse(
			http.StatusBadRequest,
			rs.getErrorListResponse(r.Context(), customError),
		)

		return // OK, пользовательская ошибка обработана
	}

	if customErrorList, ok := err.(mrerr.CustomErrors); ok { //nolint:errorlint
		for _, customError := range customErrorList {
			if !customError.IsValid() {
				err = customError.Err() // берётся первая попавшаяся необработанная ошибка

				goto InternalErrorSection
			}
		}

		sendResponse(
			http.StatusBadRequest,
			rs.getErrorListResponse(r.Context(), customErrorList...),
		)

		return // OK, пользовательская ошибка обработана
	}

InternalErrorSection:

	// сюда приходят следующие виды ошибок:
	// 1. AppError + Internal/System/User;
	// 2. ProtoAppError + Internal/System/User (нужно найти место их создания и добавить у для них вызов New()/Wrap());
	// 3. остальные ошибки: которые не были обёрнуты в ProtoAppError (нужно найти место их создания и обернуть);
	rs.errorHandler.PerformWithCommit(
		ctx,
		err,
		func(errType mrcore.AnalyzedErrorType, err *mrerr.AppError) {
			sendResponse(
				rs.statusGetter.ErrorStatus(errType, err),
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

		mrlog.Ctx(ctx).Error().Err(mrcore.ErrHttpResponseParseData.Wrap(err)).Msg("marshal failed")
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	mrlib.Write(ctx, w, bytes)
}

func (rs *ErrorSender) getErrorListResponse(ctx context.Context, errors ...*mrerr.CustomError) ErrorListResponse {
	attrs := make([]ErrorAttribute, len(errors))

	for i, customError := range errors {
		attrs[i].ID = customError.CustomCode()
		attrs[i].Value = customError.Err().Translate(mrlang.Ctx(ctx)).Reason

		if rs.isDebug {
			attrs[i].DebugInfo = rs.debugInfo(customError.Err())
		}
	}

	return attrs
}

func (rs *ErrorSender) getErrorDetailsResponse(r *http.Request, appError *mrerr.AppError) ErrorDetailsResponse {
	errMessage := appError.Translate(mrlang.Ctx(r.Context()))
	response := ErrorDetailsResponse{
		Title:        errMessage.Reason,
		Details:      errMessage.DetailsToString(),
		Request:      r.Method + " " + r.URL.Path,
		Time:         time.Now().UTC().Format(time.RFC3339),
		ErrorTraceID: appError.InstanceID(),
	}

	if rs.isDebug {
		if response.Details != "" {
			response.Details += "\n"
		}

		response.Details += "DebugInfo: " + rs.debugInfo(appError)
	}

	return response
}

func (rs *ErrorSender) debugInfo(err *mrerr.AppError) string {
	return fmt.Sprintf(
		"errCode=%s; errKind=%s; err={%s}",
		err.Code(),
		err.Kind(),
		err.Error(),
	)
}
