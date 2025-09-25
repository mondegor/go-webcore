package mrreq

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrdebug"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// FormFiles - возвращает список файлов из внешнего источника (multipart/form-data) по указанному ключу.
func FormFiles(r *http.Request, logger mrlog.Logger, key string, maxMemory int64) ([]*multipart.FileHeader, error) {
	if maxMemory < 1 {
		maxMemory = defaultMaxMemory
	}

	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			mrdebug.MultipartForm(r.Context(), logger, r.MultipartForm)

			if errors.Is(err, http.ErrMissingBoundary) {
				return nil, mr.ErrHttpFileUpload.Wrap(err, key)
			}

			return nil, mr.ErrHttpMultipartFormFile.Wrap(err, key)
		}
	}

	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if fhs, ok := r.MultipartForm.File[key]; ok {
			for i := range fhs {
				mrdebug.MultipartFileHeader(r.Context(), logger, fhs[i])
			}

			return fhs, nil
		}
	}

	return nil, nil
}

// FormFile - возвращает файл из внешнего источника (multipart/form-data) по указанному ключу.
func FormFile(r *http.Request, logger mrlog.Logger, key string) (*multipart.FileHeader, error) {
	fhs, err := FormFiles(r, logger, key, 0)
	if err != nil {
		return nil, err
	}

	if len(fhs) == 0 {
		return nil, mr.ErrHttpFileUpload.Wrap(http.ErrMissingFile, key)
	}

	return fhs[0], nil
}
