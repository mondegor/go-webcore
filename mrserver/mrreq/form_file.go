package mrreq

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrdebug"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// FormFiles - возвращает список файлов из внешнего источника (multipart/form-data) по указанному ключу.
func FormFiles(r *http.Request, key string, maxMemory int64) ([]*multipart.FileHeader, error) {
	if maxMemory < 1 {
		maxMemory = defaultMaxMemory
	}

	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			mrdebug.MultipartForm(r.Context(), r.MultipartForm)

			if errors.Is(err, http.ErrMissingBoundary) {
				return nil, mrcore.ErrHttpFileUpload.Wrap(err, key)
			}

			return nil, mrcore.ErrHttpMultipartFormFile.Wrap(err, key)
		}
	}

	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if fhs, ok := r.MultipartForm.File[key]; ok {
			for i := range fhs {
				mrdebug.MultipartFileHeader(r.Context(), fhs[i])
			}

			return fhs, nil
		}
	}

	return nil, nil
}

// FormFile - возвращает файл из внешнего источника (multipart/form-data) по указанному ключу.
func FormFile(r *http.Request, key string) (*multipart.FileHeader, error) {
	fhs, err := FormFiles(r, key, 0)
	if err != nil {
		return nil, err
	}

	if len(fhs) == 0 {
		return nil, mrcore.ErrHttpFileUpload.Wrap(http.ErrMissingFile, key)
	}

	return fhs[0], nil
}
