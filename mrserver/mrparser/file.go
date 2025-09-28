package mrparser

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/mondegor/go-sysmess/mrdto"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlib/extfile"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
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
		maxFiles                int
		checkRequestContentType bool
		allowedMimeTypes        *extfile.MimeTypeList
		logger                  mrlog.Logger
	}
)

// NewFile - создаёт объект File.
func NewFile(logger mrlog.Logger, opts ...FileOption) *File {
	f := &File{
		minSize:                 defaultMinSize,
		maxSize:                 defaultMaxSize,
		maxFiles:                defaultMaxFiles,
		checkRequestContentType: defaultCheckRequestContentType,
		allowedMimeTypes: extfile.NewMimeTypeList( // by default
			[]extfile.MimeType{
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
	hdr, err := mrreq.FormFile(r, p.logger, key)
	if err != nil {
		return mrtype.File{}, err
	}

	if err = p.checkFile(hdr); err != nil {
		return mrtype.File{}, err
	}

	file, err := hdr.Open()
	if err != nil {
		return mrtype.File{}, mr.ErrHttpMultipartFormFile.Wrap(err, key)
	}

	return mrtype.File{
		FileInfo: mrdto.FileInfo{
			ContentType:  p.detectedContentType(hdr),
			OriginalName: hdr.Filename,
			Size:         uint64(hdr.Size), //nolint:gosec
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

	var buf bytes.Buffer

	if _, err = buf.ReadFrom(file.Body); err != nil {
		return mrtype.FileContent{}, mr.ErrInternal.Wrap(err)
	}

	return mrtype.FileContent{
		FileInfo: file.FileInfo,
		Body:     buf.Bytes(),
	}, nil
}

// FormFiles - возвращает массив заголовков на файлы из MultipartForm.
func (p *File) FormFiles(r *http.Request, key string) ([]mrtype.FileHeader, error) {
	fds, err := mrreq.FormFiles(r, p.logger, key, 0)
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

	files := make([]mrtype.FileHeader, 0, countFiles)

	for i := 0; i < countFiles; i++ {
		if err = p.checkFile(fds[i]); err != nil {
			return nil, err
		}

		files = append(
			files,
			mrtype.FileHeader{
				FileInfo: mrdto.FileInfo{
					ContentType:  p.detectedContentType(fds[i]),
					OriginalName: fds[i].Filename,
					Size:         uint64(fds[i].Size), //nolint:gosec
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
	if hdr.Size < 0 {
		return mr.ErrValidateFileSize.New()
	}

	if uint64(hdr.Size) < p.minSize {
		return mr.ErrValidateFileSizeMin.New(p.minSize)
	}

	if p.maxSize > 0 && uint64(hdr.Size) > p.maxSize {
		return mr.ErrValidateFileSizeMax.New(p.maxSize)
	}

	detectedContentType, err := p.allowedMimeTypes.ContentTypeByExt(path.Ext(hdr.Filename))
	if err != nil {
		return mr.ErrValidateFileExtension.Wrap(err, path.Ext(hdr.Filename))
	}

	if p.checkRequestContentType {
		if detectedContentType != hdr.Header.Get(headerContentType) {
			return mr.ErrValidateFileContentType.New(hdr.Header.Get(headerContentType))
		}
	} else {
		if detectedContentType == "" {
			detectedContentType = hdr.Header.Get(headerContentType)
		}
	}

	if detectedContentType == "" {
		return mr.ErrValidateFileUnsupportedType.New(hdr.Filename)
	}

	return nil
}

func (p *File) checkTotalSize(fds []*multipart.FileHeader, countFiles int) error {
	if p.maxTotalSize > 0 {
		var currentSize uint64

		for i := 0; i < countFiles; i++ {
			if fds[i].Size < 0 {
				continue // игнорируются отрицательный размер файла, ошибка произойдёт в checkFile()
			}

			currentSize += uint64(fds[i].Size) //nolint:gosec
		}

		if currentSize > p.maxTotalSize {
			return mr.ErrValidateFileTotalSizeMax.New(p.maxTotalSize)
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
