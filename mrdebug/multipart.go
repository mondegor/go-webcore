package mrdebug

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/mondegor/go-webcore/mrlog"
)

func MultipartForm(ctx context.Context, form *multipart.Form) {
	logger := mrlog.Ctx(ctx).With().Str("func", "MultipartForm").Logger()

	if form == nil {
		logger.Debug().Msg("Param form is nil")
		return
	}

	if form.Value != nil && len(form.Value) > 0 {
		for key, values := range form.Value {
			logger.Debug().Msgf("key=%s; value=%s", key, strings.Join(values, ", "))
		}
	} else {
		logger.Debug().Msg("value is EMPTY")
	}

	if form.File != nil && len(form.File) > 0 {
		for key, fhs := range form.File {
			logger.Debug().Msgf("key=%s; fhs.len=%d", key, len(fhs))
		}
	} else {
		logger.Debug().Msg("form.File is EMPTY")
	}
}

func MultipartFileHeader(ctx context.Context, hdr *multipart.FileHeader) {
	logger := mrlog.Ctx(ctx).With().Str("func", "MultipartFileHeader").Logger()

	if hdr == nil {
		logger.Debug().Msg("Param hdr is nil")
		return
	}

	logger.Debug().Msgf(
		"uploaded file: name=%s, size=%d, header=%#v",
		hdr.Filename, hdr.Size, hdr.Header,
	)
}
