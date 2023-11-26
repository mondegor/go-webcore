package mrdebug

import (
	"mime/multipart"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

func MultipartForm(logger mrcore.Logger, form *multipart.Form) {
	if logger == nil {
		logger = mrcore.DefaultLogger()
	}

	if form == nil {
		logger.Debug("Param form *multipart.Form is nil")
		return
	}

	if form.Value != nil && len(form.Value) > 0 {
		for key, values := range form.Value {
			logger.Debug("MultipartForm.Value: key=%s; value=%s", key, strings.Join(values, ", "))
		}
	} else {
		logger.Debug("MultipartForm.Value is EMPTY")
	}

	if form.File != nil && len(form.File) > 0 {
		for key, fhs := range form.File {
			logger.Debug("MultipartForm.File: key=%s; fhs.len=%d", key, len(fhs))
		}
	} else {
		logger.Debug("MultipartForm.File is EMPTY")
	}
}

func MultipartFileHeader(logger mrcore.Logger, hdr *multipart.FileHeader) {
	if logger == nil {
		logger = mrcore.DefaultLogger()
	}

	if hdr == nil {
		logger.Debug("Param hdr *multipart.FileHeader is nil")
		return
	}

	logger.Debug(
		"uploaded file: name=%s, size=%d, header=%#v",
		hdr.Filename, hdr.Size, hdr.Header,
	)
}
