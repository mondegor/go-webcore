package mrreq

import (
	"mime/multipart"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrdebug"
)

// FormFile - WARNING you don't forget to call file.Close()
func FormFile(r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	file, hdr, err := r.FormFile(key)

	if err != nil {
		mrdebug.MultipartForm(r.Context(), r.MultipartForm)
		return nil, nil, mrcore.FactoryErrHttpMultipartFormFile.Wrap(err, key)
	}

	mrdebug.MultipartFileHeader(r.Context(), hdr)

	return file, hdr, nil
}