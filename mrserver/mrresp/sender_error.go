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
	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	ErrorSender struct {
		encoder      mrserver.ResponseEncoder
		overrideFunc mrserver.HttpErrorOverrideFunc
	}
)

var (
	// Make sure the ErrorSender conforms with the mrserver.ErrorResponseSender interface
	_ mrserver.ErrorResponseSender = (*ErrorSender)(nil)
)

func NewErrorSender(encoder mrserver.ResponseEncoder) *ErrorSender {
	return &ErrorSender{
		encoder:      encoder,
		overrideFunc: mrserver.DefaultHttpErrorOverrideFunc,
	}
}

func NewErrorSenderWithOverrideFunc(
	encoder mrserver.ResponseEncoder,
	overrideFunc mrserver.HttpErrorOverrideFunc,
) *ErrorSender {
	return &ErrorSender{
		encoder:      encoder,
		overrideFunc: overrideFunc,
	}
}

func (rs *ErrorSender) SendError(w http.ResponseWriter, r *http.Request, err error) {
	if customError, ok := err.(*mrerr.CustomError); ok {
		rs.sendStructResponse(
			r.Context(),
			w,
			http.StatusBadRequest,
			rs.encoder.ContentType(),
			rs.getErrorListResponse(r.Context(), customError),
		)

		return
	}

	if customErrorList, ok := err.(mrerr.CustomErrorList); ok {
		rs.sendStructResponse(
			r.Context(),
			w,
			http.StatusBadRequest,
			rs.encoder.ContentType(),
			rs.getErrorListResponse(r.Context(), customErrorList...),
		)

		return
	}

	appError, ok := err.(*mrerr.AppError)

	if !ok {
		appError = mrcore.FactoryErrInternalNotice.Wrap(err)
	}

	status, appError := rs.overrideFunc(appError)

	if appError.HasCallStack() {
		mrlog.Ctx(r.Context()).Error().Err(appError).Send()
	}

	rs.sendStructResponse(
		r.Context(),
		w,
		status,
		rs.encoder.ContentTypeProblem(),
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
		status = http.StatusTeapot
		bytes = []byte{}
		mrlog.Ctx(ctx).
			Error().
			Err(mrcore.FactoryErrHttpResponseParseData.Wrap(err)).
			Send()
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)

	if len(bytes) < 1 {
		return
	}

	w.Write(bytes)
}

func (rs *ErrorSender) getErrorListResponse(ctx context.Context, errors ...*mrerr.CustomError) ErrorListResponse {
	attrs := make([]ErrorAttribute, len(errors))

	for i, customError := range errors {
		attrs[i].ID = customError.Code()
		attrs[i].Value = customError.AppError().Translate(mrlang.Ctx(ctx)).Reason

		if mrdebug.IsDebug() {
			attrs[i].DebugInfo = rs.debugInfo(customError.AppError())
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

	if mrdebug.IsDebug() {
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
