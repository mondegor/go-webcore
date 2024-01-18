package mrreq

import (
	"bytes"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrtype"
)

// File - WARNING you don't forget to call result.Body.Close()
func File(r *http.Request, key string) (mrtype.File, error) {
	logger := mrctx.Logger(r.Context())
	file, hdr, err := r.FormFile(key)

	if err != nil {
		mrdebug.MultipartForm(logger, r.MultipartForm)
		return mrtype.File{}, mrcore.FactoryErrHttpMultipartFormFile.Wrap(err, key)
	}

	mrdebug.MultipartFileHeader(logger, hdr)

	return mrtype.File{
		FileInfo: mrtype.FileInfo{
			ContentType:  hdr.Header.Get("Content-Type"),
			OriginalName: hdr.Filename,
			Size:         hdr.Size,
		},
		Body: file,
	}, nil
}

// FileContent - only for short files
func FileContent(r *http.Request, key string) (mrtype.FileContent, error) {
	file, err := File(r, key)

	if err != nil {
		return mrtype.FileContent{}, nil
	}

	defer file.Body.Close()

	buf := bytes.Buffer{}

	if _, err = buf.ReadFrom(file.Body); err != nil {
		return mrtype.FileContent{}, mrcore.FactoryErrInternal.Wrap(err)
	}

	return mrtype.FileContent{
		FileInfo: file.FileInfo,
		Body:     buf.Bytes(),
	}, nil
}
