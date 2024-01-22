package mrresponse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	ErrorSender struct {
		encoder      mrserver.ResponseEncoder
		overrideFunc mrserver.HttpErrorOverrideFunc
	}
)

// Make sure the ErrorSender conforms with the mrserver.ErrorResponseSender interface
var _ mrserver.ErrorResponseSender = (*ErrorSender)(nil)

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
		mrctx.Logger(ctx).Err(mrcore.FactoryErrHttpResponseParseData.Caller(1).Wrap(err))
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
		attrs[i].Value = customError.AppError().Translate(mrctx.Locale(ctx)).Reason

		if mrcore.Debug() {
			attrs[i].DebugInfo = rs.debugInfo(customError.AppError())
		}
	}

	return attrs
}

func (rs *ErrorSender) getErrorDetailsResponse(r *http.Request, appError *mrerr.AppError) ErrorDetailsResponse {
	errMessage := appError.Translate(mrctx.Locale(r.Context()))
	response := ErrorDetailsResponse{
		Title:   errMessage.Reason,
		Details: errMessage.DetailsToString(),
		Request: r.URL.Path,
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	if mrcore.Debug() {
		if response.Details != "" {
			response.Details += "\n"
		}

		response.Details += "DebugInfo: " + rs.debugInfo(appError)
	}

	if appError.Kind() != mrerr.ErrorKindUser {
		response.ErrorTraceID = rs.getErrorTraceID(r.Context(), appError)
		mrctx.Logger(r.Context()).Err(appError)
	}

	return response
}

func (rs *ErrorSender) getErrorTraceID(ctx context.Context, err *mrerr.AppError) string {
	errorTraceID := err.TraceID()

	if errorTraceID == "" {
		return mrctx.CorrelationID(ctx)
	}

	return mrctx.CorrelationID(ctx) + ", " + errorTraceID
}

func (rs *ErrorSender) debugInfo(err *mrerr.AppError) string {
	return fmt.Sprintf(
		"errCode=%s; errKind=%s; err={%s}",
		err.Code(),
		err.Kind(),
		err.Error(),
	)
}
