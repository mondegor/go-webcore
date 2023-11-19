package mrserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	clientContext struct {
		request          *http.Request
		responseWriter   http.ResponseWriter
		pathParams       httprouter.Params
		errorWrapperFunc mrcore.ClientErrorWrapperFunc
		tools            mrctx.ClientTools
	}
)

// Make sure the clientContext conforms with the mrcore.ClientContext interface
var _ mrcore.ClientContext = (*clientContext)(nil)

func newClientData(
	r *http.Request,
	w http.ResponseWriter,
	ew mrcore.ClientErrorWrapperFunc,
	tools mrctx.ClientTools,
) *clientContext {
	if ew == nil {
		ew = DefaultErrorWrapperFunc
	}

	c := clientContext{
		request:          r,
		responseWriter:   w,
		errorWrapperFunc: ew,
		tools:            tools,
	}

	params, ok := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)

	if ok {
		c.pathParams = params
	}

	return &c
}

func (c *clientContext) Request() *http.Request {
	return c.request
}

func (c *clientContext) ParamFromPath(name string) string {
	return c.pathParams.ByName(name)
}

func (c *clientContext) Context() context.Context {
	return c.request.Context()
}

func (c *clientContext) WithContext(ctx context.Context) mrcore.ClientContext {
	return &clientContext{
		request:          c.request.WithContext(ctx),
		responseWriter:   c.responseWriter,
		pathParams:       c.pathParams,
		errorWrapperFunc: c.errorWrapperFunc,
		tools:            c.tools,
	}
}

func (c *clientContext) Writer() http.ResponseWriter {
	return c.responseWriter
}

func (c *clientContext) Validate(structRequest any) error {
	if err := c.parse(structRequest); err != nil {
		return err
	}

	return c.validate(structRequest)
}

func (c *clientContext) parse(structRequest any) error {
	dec := json.NewDecoder(c.request.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&structRequest); err != nil {
		return mrcore.FactoryErrHttpRequestParseData.Wrap(err)
	}

	return nil
}

func (c *clientContext) validate(structRequest any) error {
	return c.tools.Validator.Validate(c.request.Context(), structRequest)
}

func (c *clientContext) SendResponse(status int, structResponse any) error {
	// WARNING: err is important, because error(nil) != *mrerr.AppError(nil)
	if err := c.sendResponse(status, structResponse); err != nil {
		return err
	}

	return nil
}

func (c *clientContext) sendResponse(status int, structResponse any) *mrerr.AppError {
	c.responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.responseWriter.WriteHeader(status)

	bytes, err := json.Marshal(structResponse)

	if err != nil {
		return mrcore.FactoryErrHttpResponseParseData.Wrap(err)
	}

	_, err = c.responseWriter.Write(bytes)

	if err != nil {
		return mrcore.FactoryErrHttpResponseSendData.Wrap(err)
	}

	return nil
}

func (c *clientContext) SendResponseNoContent() error {
	c.responseWriter.WriteHeader(http.StatusNoContent)

	return nil
}

func (c *clientContext) SendFile(info mrtype.FileInfo, attachmentName string, file io.Reader) error {
	c.responseWriter.Header().Set("Content-Type", info.ContentType)

	if info.Size > 0 {
		c.responseWriter.Header().Set("Content-Length", strconv.FormatInt(info.Size, 10))
	}

	if attachmentName != "" {
		c.responseWriter.Header().Set("Cache-control", "private")
		c.responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachmentName)) // :TODO: escape
	}

	c.responseWriter.WriteHeader(http.StatusOK)

	_, err := io.Copy(c.responseWriter, file)

	if err != nil {
		return err
	}

	return nil
}

func (c *clientContext) sendErrorResponse(err error) {
	var appError *mrerr.AppError

	for { // only for break
		if fieldError, ok := err.(*mrerr.FieldError); ok {
			appError = c.sendFieldErrorListResponse(mrerr.FieldErrorList{fieldError})
			break
		}

		if fieldErrorList, ok := err.(mrerr.FieldErrorList); ok {
			appError = c.sendFieldErrorListResponse(fieldErrorList)
			break
		}

		if appErrorTmp, ok := err.(*mrerr.AppError); ok {
			if appErrorTmp.Kind() == mrerr.ErrorKindUser {
				appError = c.sendSystemErrorResponse(appErrorTmp)
				break
			}

			appError = appErrorTmp
			break
		}

		appError = mrcore.FactoryErrInternal.Caller(-1).Wrap(err)
		break
	}

	if appError != nil {
		c.tools.Logger.DisableFileLine().Err(appError)
		c.sendAppErrorResponse(c.errorWrapperFunc(appError))
	}
}

func (c *clientContext) sendFieldErrorListResponse(list mrerr.FieldErrorList) *mrerr.AppError {
	errorResponseList := AppErrorListResponse{}

	for _, fieldError := range list {
		if fieldError.Kind() != mrerr.ErrorKindUser {
			c.tools.Logger.Err(fieldError.AppError())
			continue
		}

		errorResponseList.Add(
			fieldError.ID(),
			fieldError.AppError().Translate(c.tools.Locale).Reason,
		)
	}

	return c.sendResponse(http.StatusBadRequest, errorResponseList)
}

func (c *clientContext) sendSystemErrorResponse(appError *mrerr.AppError) *mrerr.AppError {
	return c.sendResponse(
		http.StatusBadRequest,
		AppErrorListResponse{
			AppErrorAttribute{
				ID:    AppErrorAttributeNameSystem,
				Value: appError.Translate(c.tools.Locale).Reason,
			},
		},
	)
}

func (c *clientContext) sendAppErrorResponse(status int, appError *mrerr.AppError) {
	c.responseWriter.Header().Set("Content-Type", "application/problem+json")
	c.responseWriter.WriteHeader(status)

	errMessage := appError.Translate(c.tools.Locale)
	errorResponse := AppErrorResponse{
		Title:        errMessage.Reason,
		Details:      errMessage.DetailsToString(),
		Request:      c.request.URL.Path,
		Time:         time.Now().Format(time.RFC3339),
		ErrorTraceID: c.getErrorTraceID(appError),
	}

	c.responseWriter.Write(errorResponse.Marshal())
}

func (c *clientContext) getErrorTraceID(err *mrerr.AppError) string {
	errorTraceID := err.TraceID()

	if errorTraceID == "" {
		return c.tools.CorrelationID
	}

	return c.tools.CorrelationID + ", " + errorTraceID
}
