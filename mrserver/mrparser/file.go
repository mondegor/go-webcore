package mrparser

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// File - парсер файлов.
	File struct {
		allowedMimeTypes        *mrlib.MimeTypeList
		minSize                 int64
		maxSize                 int64
		maxTotalSize            int64
		maxFiles                int
		checkRequestContentType bool
	}

	// FileOptions - опции для создания объекта File.
	FileOptions struct {
		AllowedMimeTypes        *mrlib.MimeTypeList
		MinSize                 int64
		MaxSize                 int64
		MaxTotalSize            int64
		MaxFiles                int
		CheckRequestContentType bool
	}
)

// Make sure the File conforms with the mrserver.RequestParserFile interface.
var _ mrserver.RequestParserFile = (*File)(nil)

// NewFile - создаёт объект File.
func NewFile(logger mrlog.Logger, opts FileOptions) *File {
	if opts.AllowedMimeTypes == nil {
		opts.AllowedMimeTypes = mrlib.NewMimeTypeList(
			logger,
			[]mrlib.MimeType{
				{
					Extension:   ".pdf",
					ContentType: "application/pdf",
				},
				{
					Extension:   ".zip",
					ContentType: "application/zip",
				},
			},
		)
	}

	if opts.MaxTotalSize == 0 && opts.MaxSize > 0 && opts.MaxFiles > 0 {
		opts.MaxTotalSize = opts.MaxSize * int64(opts.MaxFiles)
	}

	return &File{
		allowedMimeTypes:        opts.AllowedMimeTypes,
		minSize:                 opts.MinSize,
		maxSize:                 opts.MaxSize,
		maxFiles:                opts.MaxFiles,
		maxTotalSize:            opts.MaxTotalSize,
		checkRequestContentType: opts.CheckRequestContentType,
	}
}

// FormFile - возвращает информацию о файле со ссылкой для чтения файла из MultipartForm.
// WARNING: you don't forget to call result.Body.Close().
func (p *File) FormFile(r *http.Request, key string) (mrtype.File, error) {
	hdr, err := mrreq.FormFile(r, key)
	if err != nil {
		return mrtype.File{}, err
	}

	if err = p.checkFile(hdr); err != nil {
		return mrtype.File{}, err
	}

	file, err := hdr.Open()
	if err != nil {
		return mrtype.File{}, mrcore.ErrHttpMultipartFormFile.Wrap(err, key)
	}

	return mrtype.File{
		FileInfo: mrtype.FileInfo{
			ContentType:  p.detectedContentType(hdr),
			OriginalName: hdr.Filename,
			Size:         hdr.Size,
		},
		Body: file,
	}, nil
}

// FormFileContent - возвращает информацию о файле и сам файл из MultipartForm.
// WARNING: only for short files.
func (p *File) FormFileContent(r *http.Request, key string) (mrtype.FileContent, error) {
	file, err := p.FormFile(r, key)
	if err != nil {
		return mrtype.FileContent{}, err
	}

	defer file.Body.Close()

	buf := bytes.Buffer{}

	if _, err = buf.ReadFrom(file.Body); err != nil {
		return mrtype.FileContent{}, mrcore.ErrInternal.Wrap(err)
	}

	return mrtype.FileContent{
		FileInfo: file.FileInfo,
		Body:     buf.Bytes(),
	}, nil
}

// FormFiles - возвращает массив заголовков на файлы из MultipartForm.
func (p *File) FormFiles(r *http.Request, key string) ([]mrtype.FileHeader, error) {
	fds, err := mrreq.FormFiles(r, key, 0)
	if err != nil {
		return nil, err
	}

	if len(fds) == 0 {
		return nil, nil
	}

	countFiles := p.allowedFiles(len(fds))

	if err = p.checkTotalSize(fds, countFiles); err != nil {
		return nil, err
	}

	files := make([]mrtype.FileHeader, 0, countFiles)

	for i := 0; i < countFiles; i++ {
		if err = p.checkFile(fds[i]); err != nil {
			return nil, err
		}

		files = append(
			files,
			mrtype.FileHeader{
				FileInfo: mrtype.FileInfo{
					ContentType:  p.detectedContentType(fds[i]),
					OriginalName: fds[i].Filename,
					Size:         fds[i].Size,
				},
				Header: fds[i],
			},
		)
	}

	return files, nil
}

func (p *File) checkFile(hdr *multipart.FileHeader) error {
	if hdr.Size < p.minSize {
		return ErrHttpRequestFileSizeMin.New(p.minSize)
	}

	if p.maxSize > 0 && hdr.Size > p.maxSize {
		return ErrHttpRequestFileSizeMax.New(p.maxSize)
	}

	ext := path.Ext(hdr.Filename)

	if err := p.allowedMimeTypes.CheckExt(ext); err != nil {
		return ErrHttpRequestFileExtension.Wrap(err, ext)
	}

	detectedContentType := p.allowedMimeTypes.ContentType(ext)

	if p.checkRequestContentType {
		if detectedContentType != hdr.Header.Get("Content-Type") {
			return ErrHttpRequestFileContentType.New(hdr.Header.Get("Content-Type"))
		}
	} else {
		if detectedContentType == "" {
			detectedContentType = hdr.Header.Get("Content-Type")
		}
	}

	if detectedContentType == "" {
		return ErrHttpRequestFileUnsupportedType.New(hdr.Filename)
	}

	return nil
}

func (p *File) allowedFiles(count int) int {
	if p.maxFiles > 0 && count > p.maxFiles {
		return p.maxFiles
	}

	return count
}

func (p *File) checkTotalSize(fds []*multipart.FileHeader, countFiles int) error {
	if p.maxTotalSize > 0 {
		var currentSize int64

		for i := 0; i < countFiles; i++ {
			currentSize += fds[i].Size
		}

		if currentSize > p.maxTotalSize {
			return ErrHttpRequestFileTotalSizeMax.New(p.maxTotalSize)
		}
	}

	return nil
}

func (p *File) detectedContentType(hdr *multipart.FileHeader) string {
	if contentType := p.allowedMimeTypes.ContentTypeByFileName(hdr.Filename); contentType != "" {
		return contentType
	}

	return hdr.Header.Get("Content-Type")
}
