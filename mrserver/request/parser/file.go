package parser

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrmodel"
	"github.com/mondegor/go-sysmess/util/mime"

	"github.com/mondegor/go-webcore/mrdebug"
)

const (
	headerContentType              = "Content-Type"
	defaultMinSize                 = 0           // bytes
	defaultMaxSize                 = 1024 * 1024 // 1Mb
	defaultMaxFiles                = 4
	defaultMaxMemory               = 32 << 20 // 32 MB
	defaultCheckRequestContentType = false
)

type (
	// File - парсер файлов.
	File struct {
		minSize                 int64 // bytes
		maxSize                 int64 // bytes
		maxTotalSize            int64 // bytes
		maxFiles                int
		checkRequestContentType bool
		allowedMimeTypes        *mime.TypeList
		logger                  mrlog.Logger
	}
)

// NewFile - создаёт объект File.
func NewFile(logger mrlog.Logger, opts ...FileOption) *File {
	o := fileOptions{
		file: &File{
			minSize:                 defaultMinSize,
			maxSize:                 defaultMaxSize,
			maxFiles:                defaultMaxFiles,
			checkRequestContentType: defaultCheckRequestContentType,
			allowedMimeTypes: mime.NewTypeList( // by default
				[]mime.Type{
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
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	// вычисление и установка maxTotalSize
	if o.file.maxSize > 0 {
		o.file.maxTotalSize = o.file.maxSize

		if o.file.maxFiles > 0 {
			o.file.maxTotalSize *= int64(o.file.maxFiles)
		}
	}

	return o.file
}

// FormFile - возвращает информацию о файле со ссылкой для чтения файла из MultipartForm.
// WARNING: you don't forget to call result.Body.Close().
func (p *File) FormFile(r *http.Request, key string) (mrmodel.File, error) {
	hdr, err := p.formFile(r, p.logger, key)
	if err != nil {
		return mrmodel.File{}, err
	}

	if err = p.checkFile(hdr); err != nil {
		return mrmodel.File{}, err
	}

	file, err := hdr.Open()
	if err != nil {
		return mrmodel.File{}, errors.ErrSystemHttpMultipartFormFile.Wrap(err, "key", key)
	}

	return mrmodel.File{
		FileInfo: mrmodel.FileInfo{
			ContentType:  p.detectedContentType(hdr),
			OriginalName: hdr.Filename,
			Size:         hdr.Size,
		},
		Body: file,
	}, nil
}

// FormFileContent - возвращает информацию о файле и сам файл из MultipartForm.
// WARNING: only for short files.
func (p *File) FormFileContent(r *http.Request, key string) (mrmodel.FileContent, error) {
	file, err := p.FormFile(r, key)
	if err != nil {
		return mrmodel.FileContent{}, err
	}

	defer func() {
		_ = file.Body.Close()
	}()

	var buf bytes.Buffer

	if _, err = buf.ReadFrom(file.Body); err != nil {
		return mrmodel.FileContent{}, errors.WrapInternalError(err, "reading file.Body failed")
	}

	return mrmodel.FileContent{
		FileInfo: file.FileInfo,
		Body:     buf.Bytes(),
	}, nil
}

// FormFiles - возвращает массив заголовков на файлы из MultipartForm.
func (p *File) FormFiles(r *http.Request, key string) ([]mrmodel.FileHeader, error) {
	fds, err := p.formFiles(r, p.logger, key, 0)
	if err != nil {
		return nil, err
	}

	if len(fds) == 0 {
		return nil, nil
	}

	countFiles := len(fds)

	if p.maxFiles > 0 && countFiles > p.maxFiles {
		countFiles = p.maxFiles
	}

	if err = p.checkTotalSize(fds, countFiles); err != nil {
		return nil, err
	}

	files := make([]mrmodel.FileHeader, 0, countFiles)

	for i := 0; i < countFiles; i++ {
		if err = p.checkFile(fds[i]); err != nil {
			return nil, err
		}

		files = append(
			files,
			mrmodel.FileHeader{
				FileInfo: mrmodel.FileInfo{
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
	if hdr.Size < 0 {
		return errors.ErrValidateFileSize
	}

	if hdr.Size < p.minSize {
		return errors.ErrValidateFileSizeMin.New(p.minSize)
	}

	if p.maxSize > 0 && hdr.Size > p.maxSize {
		return errors.ErrValidateFileSizeMax.New(p.maxSize)
	}

	detectedContentType, err := p.allowedMimeTypes.ContentTypeByExt(path.Ext(hdr.Filename))
	if err != nil {
		return errors.ErrValidateFileExtension.Wrap(err, path.Ext(hdr.Filename))
	}

	if p.checkRequestContentType {
		if detectedContentType != hdr.Header.Get(headerContentType) {
			return errors.ErrValidateFileContentType.New(hdr.Header.Get(headerContentType))
		}
	} else {
		if detectedContentType == "" {
			detectedContentType = hdr.Header.Get(headerContentType)
		}
	}

	if detectedContentType == "" {
		return errors.ErrValidateFileUnsupportedType.New(hdr.Filename)
	}

	return nil
}

func (p *File) checkTotalSize(fds []*multipart.FileHeader, countFiles int) error {
	if p.maxTotalSize > 0 {
		var currentSize int64

		for i := 0; i < countFiles; i++ {
			if fds[i].Size < 0 {
				continue // игнорируются отрицательный размер файла, ошибка произойдёт в checkFile()
			}

			currentSize += fds[i].Size
		}

		if currentSize > p.maxTotalSize {
			return errors.ErrValidateFileTotalSizeMax.New(p.maxTotalSize)
		}
	}

	return nil
}

func (p *File) detectedContentType(hdr *multipart.FileHeader) string {
	contentType, err := p.allowedMimeTypes.ContentTypeByExt(path.Ext(hdr.Filename))
	if err != nil || contentType == "" {
		return hdr.Header.Get(headerContentType)
	}

	return contentType
}

// formFiles - возвращает список файлов из внешнего источника (multipart/form-data) по указанному ключу.
func (p *File) formFiles(r *http.Request, logger mrlog.Logger, key string, maxMemory int64) ([]*multipart.FileHeader, error) {
	if maxMemory < 1 {
		maxMemory = defaultMaxMemory
	}

	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			mrdebug.MultipartForm(r.Context(), logger, r.MultipartForm)

			if errors.Is(err, http.ErrMissingBoundary) {
				return nil, errors.ErrHttpFileUpload.Wrap(err, key)
			}

			return nil, errors.ErrSystemHttpMultipartFormFile.Wrap(err, "key", key)
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

// formFile - возвращает файл из внешнего источника (multipart/form-data) по указанному ключу.
func (p *File) formFile(r *http.Request, logger mrlog.Logger, key string) (*multipart.FileHeader, error) {
	fhs, err := p.formFiles(r, logger, key, 0)
	if err != nil {
		return nil, err
	}

	if len(fhs) == 0 {
		return nil, errors.ErrHttpFileUpload.Wrap(http.ErrMissingFile, key)
	}

	return fhs[0], nil
}
