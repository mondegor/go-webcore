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

const (
	contentTypeJson        = "application/json; charset=utf-8"
	contentTypeProblemJson = "application/problem+json; charset=utf-8"
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
		c.tools.Logger.Caller(1).Warn(err)
		return mrcore.FactoryErrHttpRequestParseData.Wrap(err)
	}

	return nil
}

func (c *clientContext) validate(structRequest any) error {
	return c.tools.Validator.Validate(c.request.Context(), structRequest)
}

func (c *clientContext) SendResponse(status int, structResponse any) error {
	bytes, err := json.Marshal(structResponse)

	if err != nil {
		return mrcore.FactoryErrHttpResponseParseData.Wrap(err)
	}

	c.sendResponse(status, contentTypeJson, bytes)

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

func (c *clientContext) sendResponse(status int, contentType string, body []byte) {
	c.responseWriter.Header().Set("Content-Type", contentType)
	c.responseWriter.WriteHeader(status)

	if len(body) < 1 {
		return
	}

	_, err := c.responseWriter.Write(body)

	if err != nil {
		c.tools.Logger.DisableFileLine().Err(mrcore.FactoryErrHttpResponseSendData.Caller(1).Wrap(err))
	}
}

func (c *clientContext) sendStructResponse(status int, contentType string, structResponse any) {
	bytes, err := json.Marshal(structResponse)

	if err != nil {
		status = http.StatusTeapot
		bytes = []byte{}
		c.tools.Logger.DisableFileLine().Err(mrcore.FactoryErrHttpResponseParseData.Caller(1).Wrap(err))
	}

	c.sendResponse(status, contentType, bytes)
}

func (c *clientContext) sendErrorResponse(err error) {
	if fieldError, ok := err.(*mrerr.FieldError); ok {
		c.sendStructResponse(
			http.StatusBadRequest,
			contentTypeJson,
			c.getErrorListResponse(fieldError),
		)

		return
	}

	if fieldErrorList, ok := err.(mrerr.FieldErrorList); ok {
		c.sendStructResponse(
			http.StatusBadRequest,
			contentTypeJson,
			c.getErrorListResponse(fieldErrorList...),
		)

		return
	}

	appError, ok := err.(*mrerr.AppError)

	if ok {
		if appError.Kind() == mrerr.ErrorKindUser {
			c.sendStructResponse(
				http.StatusBadRequest,
				contentTypeJson,
				c.getErrorListResponse(
					mrerr.NewFieldErrorAppError(ErrorAttributeNameByDefault, appError),
				),
			)

			return
		}
	} else {
		appError = mrcore.FactoryErrInternal.Caller(-1).Wrap(err)
	}

	c.tools.Logger.DisableFileLine().Err(appError)
	status, appError := c.errorWrapperFunc(appError)

	c.sendStructResponse(
		status,
		contentTypeProblemJson,
		c.getErrorDetailsResponse(appError),
	)
}

func (c *clientContext) getErrorListResponse(fields ...*mrerr.FieldError) ErrorListResponse {
	attrs := make([]ErrorAttribute, len(fields))

	for i, fieldError := range fields {
		attrs[i].ID = fieldError.ID()
		attrs[i].Value = fieldError.AppError().Translate(c.tools.Locale).Reason

		if mrcore.Debug() {
			attrs[i].DebugInfo = c.debugInfo(fieldError.AppError())
		}
	}

	return attrs
}

func (c *clientContext) getErrorDetailsResponse(appError *mrerr.AppError) ErrorDetailsResponse {
	errMessage := appError.Translate(c.tools.Locale)
	response := ErrorDetailsResponse{
		Title:        errMessage.Reason,
		Details:      errMessage.DetailsToString(),
		Request:      c.request.URL.Path,
		Time:         time.Now().Format(time.RFC3339),
		ErrorTraceID: c.getErrorTraceID(appError),
	}

	if mrcore.Debug() {
		if response.Details != "" {
			response.Details += "\n"
		}

		response.Details += "DebugInfo: " + c.debugInfo(appError)
	}

	return response
}

func (c *clientContext) getErrorTraceID(err *mrerr.AppError) string {
	errorTraceID := err.TraceID()

	if errorTraceID == "" {
		return c.tools.CorrelationID
	}

	return c.tools.CorrelationID + ", " + errorTraceID
}

func (c *clientContext) debugInfo(err *mrerr.AppError) string {
	return fmt.Sprintf(
		"errId=%s; errKind=%s; err={%s}",
		err.ID(),
		err.Kind(),
		err.Error(),
	)
}
