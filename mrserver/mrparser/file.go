package mrparser

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	File struct {
		allowedExts             []string
		minSize                 int64
		maxSize                 int64
		checkRequestContentType bool
	}

	FileOptions struct {
		AllowedExts             []string
		MinSize                 int64
		MaxSize                 int64
		CheckRequestContentType bool
	}

	rawFile struct {
		file        multipart.File
		hdr         *multipart.FileHeader
		contentType string
	}
)

// Make sure the File conforms with the mrserver.RequestParserFile interface
var _ mrserver.RequestParserFile = (*File)(nil)

func NewFile(opts FileOptions) *File {
	if len(opts.AllowedExts) == 0 {
		opts.AllowedExts = []string{".pdf", ".json", ".rar", ".tar", ".tgz", ".zip"}
	}

	return &File{
		allowedExts:             opts.AllowedExts,
		minSize:                 opts.MinSize,
		maxSize:                 opts.MaxSize,
		checkRequestContentType: opts.CheckRequestContentType,
	}
}

// FormFile - WARNING you don't forget to call result.Body.Close()
func (p *File) FormFile(r *http.Request, key string) (mrtype.File, error) {
	raw, err := p.raw(r, key)

	if err != nil {
		return mrtype.File{}, err
	}

	return mrtype.File{
		FileInfo: mrtype.FileInfo{
			ContentType:  raw.contentType,
			OriginalName: raw.hdr.Filename,
			Size:         raw.hdr.Size,
		},
		Body: raw.file,
	}, nil
}

// FormFileContent - only for short files
func (p *File) FormFileContent(r *http.Request, key string) (mrtype.FileContent, error) {
	file, err := p.FormFile(r, key)

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

func (p *File) raw(r *http.Request, key string) (rawFile, error) {
	file, hdr, err := mrreq.FormFile(r, key)

	if err != nil {
		return rawFile{}, err
	}

	if hdr.Size < p.minSize {
		file.Close()
		return rawFile{}, FactoryErrHttpRequestFileSizeMin.New(p.minSize)
	}

	if p.maxSize > 0 && hdr.Size > p.maxSize {
		file.Close()
		return rawFile{}, FactoryErrHttpRequestFileSizeMax.New(p.maxSize)
	}

	ext := path.Ext(hdr.Filename)

	if !p.checkExt(ext) {
		file.Close()
		return rawFile{}, FactoryErrHttpRequestFileExtension.New(ext)
	}

	detectedContentType := mrlib.MimeTypeByExt(ext)

	if p.checkRequestContentType {
		if detectedContentType != hdr.Header.Get("Content-Type") {
			file.Close()
			return rawFile{}, FactoryErrHttpRequestFileContentType.New(hdr.Header.Get("Content-Type"))
		}
	} else {
		if detectedContentType == "" {
			detectedContentType = hdr.Header.Get("Content-Type")
		}
	}

	if detectedContentType == "" {
		return rawFile{}, FactoryErrHttpRequestFileUnsupportedType.New(hdr.Filename)
	}

	return rawFile{file, hdr, detectedContentType}, err
}

// ext can be empty
func (p *File) checkExt(ext string) bool {
	if len(p.allowedExts) == 0 {
		return true
	}

	for i := range p.allowedExts {
		if ext == p.allowedExts[i] {
			return true
		}
	}

	return false
}
