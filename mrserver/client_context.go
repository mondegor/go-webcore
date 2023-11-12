package mrserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	clientContext struct {
		request        *http.Request
		responseWriter http.ResponseWriter
		requestPath    *requestPath
		validator      mrcore.Validator
	}
)

// Make sure the clientContext conforms with the mrcore.ClientData interface
var _ mrcore.ClientData = (*clientContext)(nil)

func (c *clientContext) Request() *http.Request {
	return c.request
}

func (c *clientContext) RequestPath() mrcore.RequestPath {
	if c.requestPath == nil {
		c.requestPath = newRequestPath(c.request)
	}

	return c.requestPath
}

func (c *clientContext) Context() context.Context {
	return c.request.Context()
}

func (c *clientContext) WithContext(ctx context.Context) mrcore.ClientData {
	c.request = c.request.WithContext(ctx)

	return c
}

func (c *clientContext) Writer() http.ResponseWriter {
	return c.responseWriter
}

func (c *clientContext) Parse(structRequest any) error {
	dec := json.NewDecoder(c.request.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&structRequest); err != nil {
		return mrcore.FactoryErrHttpRequestParseData.Wrap(err)
	}

	return nil
}

func (c *clientContext) Validate(structRequest any) error {
	return c.validator.Validate(c.request.Context(), structRequest)
}

func (c *clientContext) ParseAndValidate(structRequest any) error {
	if err := c.Parse(structRequest); err != nil {
		return err
	}

	return c.Validate(structRequest)
}

func (c *clientContext) SendResponse(status int, structResponse any) error {
	appError := c.sendResponse(status, structResponse)

	if appError != nil {
		return appError
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
		c.responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", attachmentName)) // :TODO: escape
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
		if userErrorList, ok := err.(*mrerr.FieldErrorList); ok {
			appError = c.sendUserErrorListResponse(userErrorList)
			break
		}

		if appErrorTmp, ok := err.(*mrerr.AppError); ok {
			if appErrorTmp.Kind() == mrerr.ErrorKindUser {
				appError = c.sendUserErrorResponse(appErrorTmp)
				break
			}

			appError = appErrorTmp
			break
		}

		appError = mrcore.FactoryErrInternal.Wrap(err)
		break
	}

	if appError != nil {
		mrctx.Logger(c.Context()).DisableFileLine().Err(appError)
		c.sendAppErrorResponse(c.wrapErrorFunc(appError))
	}
}

func (c *clientContext) sendUserErrorListResponse(list *mrerr.FieldErrorList) *mrerr.AppError {
	locale := mrctx.Locale(c.Context())
	errorResponseList := AppErrorListResponse{}

	for _, userError := range *list {
		if userError.Err.Kind() != mrerr.ErrorKindUser {
			mrctx.Logger(c.Context()).Err(userError.Err)
			continue
		}

		errorResponseList.Add(
			userError.ID,
			userError.Err.Translate(locale).Reason,
		)
	}

	return c.sendResponse(http.StatusBadRequest, errorResponseList)
}

func (c *clientContext) sendUserErrorResponse(appError *mrerr.AppError) *mrerr.AppError {
	locale := mrctx.Locale(c.Context())

	return c.sendResponse(
		http.StatusBadRequest,
		AppErrorListResponse{
			AppErrorAttribute{
				ID:    AppErrorAttributeNameSystem,
				Value: appError.Translate(locale).Reason,
			},
		},
	)
}

func (c *clientContext) sendAppErrorResponse(status int, appError *mrerr.AppError) {
	c.responseWriter.Header().Set("Content-Type", "application/problem+json")
	c.responseWriter.WriteHeader(status)

	locale := mrctx.Locale(c.Context())
	errMessage := appError.Translate(locale)
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
	correlationID := mrctx.CorrelationID(c.Context())

	if errorTraceID == "" {
		return correlationID
	}

	return fmt.Sprintf("%s, %s", correlationID, errorTraceID)
}

// :TODO: move to package internal
func (c *clientContext) wrapErrorFunc(err *mrerr.AppError) (int, *mrerr.AppError) {
	status := http.StatusInternalServerError

	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		status = http.StatusNotFound
		err = mrcore.FactoryErrHttpResourceNotFound.Wrap(err)
	} else if mrcore.FactoryErrServiceTemporarilyUnavailable.Is(err) {
		err = mrcore.FactoryErrHttpResponseSystemTemporarilyUnableToProcess.Wrap(err)
	} else if err.ID() == mrerr.ErrorInternalID {
		status = http.StatusTeapot
	}

	return status, err
}
