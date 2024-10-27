package mrparser

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	headerContentType = "Content-Type"

	defaultMinSize                 = 0           // bytes
	defaultMaxSize                 = 1024 * 1024 // 1Mb
	defaultMaxFiles                = 4
	defaultCheckRequestContentType = false
)

type (
	// File - парсер файлов.
	File struct {
		minSize                 uint64 // bytes
		maxSize                 uint64 // bytes
		maxTotalSize            uint64 // bytes
		maxFiles                uint32
		checkRequestContentType bool
		allowedMimeTypes        *mrlib.MimeTypeList
		logger                  logger
	}
)

// NewFile - создаёт объект File.
func NewFile(logger logger, opts ...FileOption) *File {
	f := &File{
		minSize:                 defaultMinSize,
		maxSize:                 defaultMaxSize,
		maxFiles:                defaultMaxFiles,
		checkRequestContentType: defaultCheckRequestContentType,
		allowedMimeTypes: mrlib.NewMimeTypeList( // by default
			logger,
			[]mrlib.MimeType{
				{
					ContentType: "application/pdf",
					Extension:   ".pdf",
				},
				{
					ContentType: "application/zip",
					Extension:   ".zip",
				},
			},
		),
		logger: logger,
	}

	f.applyOptions(opts)

	// вычисление и установка maxTotalSize
	if f.maxSize > 0 {
		f.maxTotalSize = f.maxSize

		if f.maxFiles > 0 {
			f.maxTotalSize *= uint64(f.maxFiles)
		}
	}

	return f
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
			Size:         uint64(hdr.Size),
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

	countFiles := p.allowedFiles(uint32(len(fds)))

	if err = p.checkTotalSize(fds, countFiles); err != nil {
		return nil, err
	}

	files := make([]mrtype.FileHeader, 0, countFiles)

	for i := uint32(0); i < countFiles; i++ {
		if err = p.checkFile(fds[i]); err != nil {
			return nil, err
		}

		files = append(
			files,
			mrtype.FileHeader{
				FileInfo: mrtype.FileInfo{
					ContentType:  p.detectedContentType(fds[i]),
					OriginalName: fds[i].Filename,
					Size:         uint64(fds[i].Size),
				},
				Header: fds[i],
			},
		)
	}

	return files, nil
}

func (p *File) applyOptions(opts []FileOption) {
	for _, opt := range opts {
		opt(p)
	}
}

func (p *File) checkFile(hdr *multipart.FileHeader) error {
	if hdr.Size < int64(p.minSize) {
		return ErrHttpRequestFileSizeMin.New(p.minSize)
	}

	if p.maxSize > 0 && hdr.Size > int64(p.maxSize) {
		return ErrHttpRequestFileSizeMax.New(p.maxSize)
	}

	ext := path.Ext(hdr.Filename)

	if err := p.allowedMimeTypes.CheckExt(ext); err != nil {
		return ErrHttpRequestFileExtension.Wrap(err, ext)
	}

	detectedContentType := p.allowedMimeTypes.ContentType(ext)

	if p.checkRequestContentType {
		if detectedContentType != hdr.Header.Get(headerContentType) {
			return ErrHttpRequestFileContentType.New(hdr.Header.Get(headerContentType))
		}
	} else {
		if detectedContentType == "" {
			detectedContentType = hdr.Header.Get(headerContentType)
		}
	}

	if detectedContentType == "" {
		return ErrHttpRequestFileUnsupportedType.New(hdr.Filename)
	}

	return nil
}

func (p *File) allowedFiles(count uint32) uint32 {
	if p.maxFiles > 0 && count > p.maxFiles {
		return p.maxFiles
	}

	return count
}

func (p *File) checkTotalSize(fds []*multipart.FileHeader, countFiles uint32) error {
	if p.maxTotalSize > 0 {
		var currentSize int64

		for i := uint32(0); i < countFiles; i++ {
			currentSize += fds[i].Size
		}

		if currentSize > int64(p.maxTotalSize) {
			return ErrHttpRequestFileTotalSizeMax.New(p.maxTotalSize)
		}
	}

	return nil
}

func (p *File) detectedContentType(hdr *multipart.FileHeader) string {
	if contentType := p.allowedMimeTypes.ContentTypeByFileName(hdr.Filename); contentType != "" {
		return contentType
	}

	return hdr.Header.Get(headerContentType)
}
