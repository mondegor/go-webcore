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
	// ErrorSender - comment struct.
	ErrorSender struct {
		encoder          mrserver.ResponseEncoder
		errorHandler     mrcore.ErrorHandler
		statusGetter     mrserver.ErrorStatusGetter
		unexpectedStatus int
		isDebug          bool
	}
)

// Make sure the ErrorSender conforms with the mrserver.ErrorResponseSender interface.
var _ mrserver.ErrorResponseSender = (*ErrorSender)(nil)

// NewErrorSender - создаёт объект ErrorSender.
func NewErrorSender(
	encoder mrserver.ResponseEncoder,
	errorHandler mrcore.ErrorHandler,
	statusGetter mrserver.ErrorStatusGetter,
	unexpectedStatus int,
	isDebug bool,
) *ErrorSender {
	return &ErrorSender{
		encoder:          encoder,
		errorHandler:     errorHandler,
		statusGetter:     statusGetter,
		unexpectedStatus: unexpectedStatus,
		isDebug:          isDebug,
	}
}

// SendError - отправляет клиенту ответ с ошибкой в одном из статусов: 4xx, 5XX и её деталями.
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

	var appError *mrerr.AppError

	if customError, ok := err.(*mrerr.CustomError); ok { //nolint:errorlint
		if customError.IsValid() {
			sendResponse(
				http.StatusBadRequest,
				rs.getErrorListResponse(r.Context(), customError),
			)

			return
		}

		appError = customError.Err()
	}

	if customErrorList, ok := err.(mrerr.CustomErrors); ok { //nolint:errorlint
		for _, customError := range customErrorList {
			if !customError.IsValid() {
				appError = customError.Err() // берётся первая попавшаяся необработанная ошибка

				break
			}
		}

		if appError == nil {
			sendResponse(
				http.StatusBadRequest,
				rs.getErrorListResponse(r.Context(), customErrorList...),
			)

			return
		}
	}

	// сюда могут приходит следующие типы ошибок:
	// 1. AppError + Internal/System/User;
	// 1. ProtoAppError + Internal/System/User;
	// 3. error - ошибка, которая не была обёрнута в AppError;
	rs.errorHandler.Process(ctx, err)

	appError = mrcore.CastToAppError(err)

	sendResponse(
		rs.statusGetter.ErrorStatus(appError),
		rs.getErrorDetailsResponse(r, appError),
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
		status = rs.unexpectedStatus
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
		Request:      r.URL.Path,
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
