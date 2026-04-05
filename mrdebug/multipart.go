package mrdebug

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/mondegor/go-sysmess/mrlog"
)

// MultipartForm - выводит в журнал отладочную информацию о multipart-форме.
func MultipartForm(ctx context.Context, logger mrlog.Logger, form *multipart.Form) {
	if !mrlog.DebugEnabled(logger) {
		return
	}

	logger = mrlog.WithAttrs(logger, "func", "MultipartForm")

	if form == nil {
		logger.Debug(ctx, "Param form is nil")

		return
	}

	if len(form.Value) > 0 {
		for key, values := range form.Value {
			logger.Debug(ctx, fmt.Sprintf("key='%s', value='%s'", key, strings.Join(values, ", ")))
		}
	} else {
		logger.Debug(ctx, "value is EMPTY")
	}

	if len(form.File) > 0 {
		for key, fhs := range form.File {
			logger.Debug(ctx, fmt.Sprintf("key='%s', fhs.len=%d", key, len(fhs)))
		}
	} else {
		logger.Debug(ctx, "form.File is EMPTY")
	}
}

// MultipartFileHeader - выводит в журнал отладочную информацию о заголовке multipart-файла.
func MultipartFileHeader(ctx context.Context, logger mrlog.Logger, hdr *multipart.FileHeader) {
	if !mrlog.DebugEnabled(logger) {
		return
	}

	logger = mrlog.WithAttrs(logger, "func", "MultipartFileHeader")

	if hdr == nil {
		logger.Debug(ctx, "Param hdr is nil")

		return
	}

	logger.Debug(
		ctx,
		fmt.Sprintf(
			"uploaded file: name=%s, size=%d, header=%#v",
			hdr.Filename, hdr.Size, hdr.Header,
		),
	)
}
