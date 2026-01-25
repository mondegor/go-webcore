package parser

import (
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/util/mime"
	"github.com/mondegor/go-sysmess/util/ximage"
)

const (
	defaultMaxWidth  = 3840 // pixels
	defaultMaxHeight = 2160 // pixels
	defaultCheckBody = false
)

type (
	// Image - парсер изображений.
	Image struct {
		maxWidth  uint64 // pixels
		maxHeight uint64 // pixels
		checkBody bool
		file      *File
	}

	imageMeta struct {
		width  uint64
		height uint64
	}
)

// NewImage - создаёт объект Image.
func NewImage(logger mrlog.Logger, opts ...ImageOption) *Image {
	o := imageOptions{
		image: &Image{
			maxWidth:  defaultMaxWidth,
			maxHeight: defaultMaxHeight,
			checkBody: defaultCheckBody,
		},
		fileOptions: []FileOption{
			WithFileAllowedMimeTypes( // by default
				[]mime.Type{
					{
						ContentType: "image/gif",
						Extension:   ".gif",
					},
					{
						ContentType: "image/jpeg",
						Extension:   ".jpg",
					},
					{
						ContentType: "image/png",
						Extension:   ".png",
					},
				},
			),
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	o.image.file = NewFile(logger, o.fileOptions...)

	return o.image
}

// FormImage - возвращает информацию об изображении со ссылкой для чтения файла изображения из MultipartForm.
// WARNING: you don't forget to call result.Body.Close().
func (p *Image) FormImage(r *http.Request, key string) (mrtype.Image, error) {
	hdr, err := p.file.formFile(r, p.file.logger, key)
	if err != nil {
		return mrtype.Image{}, err
	}

	if err = p.file.checkFile(hdr); err != nil {
		return mrtype.Image{}, err
	}

	file, err := hdr.Open()
	if err != nil {
		return mrtype.Image{}, errors.ErrSystemHttpMultipartFormFile.Wrap(err, "key", key)
	}

	contentType := p.file.detectedContentType(hdr)

	meta, err := p.decode(file, contentType)
	if err != nil {
		if err := file.Close(); err != nil {
			p.file.logger.Error(r.Context(), "mrparser.FormImage: error when closing file image", "error", err)
		}

		return mrtype.Image{}, err
	}

	return mrtype.Image{
		ImageInfo: mrtype.ImageInfo{
			ContentType:  contentType,
			OriginalName: hdr.Filename,
			Width:        meta.width,
			Height:       meta.height,
			Size:         uint64(hdr.Size), //nolint:gosec
		},
		Body: file,
	}, nil
}

// FormImageContent - возвращает информацию об изображении и сам файл изображения из MultipartForm.
// WARNING: only for short files.
func (p *Image) FormImageContent(r *http.Request, key string) (mrtype.ImageContent, error) {
	file, err := p.FormImage(r, key)
	if err != nil {
		return mrtype.ImageContent{}, err
	}

	defer file.Body.Close()

	var buf bytes.Buffer

	if _, err = buf.ReadFrom(file.Body); err != nil {
		return mrtype.ImageContent{}, errors.WrapInternalError(err, "reading file.Body failed")
	}

	return mrtype.ImageContent{
		ImageInfo: file.ImageInfo,
		Body:      buf.Bytes(),
	}, nil
}

// FormImages - возвращает массив заголовков на файлы изображений из MultipartForm.
func (p *Image) FormImages(r *http.Request, key string) ([]mrtype.ImageHeader, error) {
	fds, err := p.file.formFiles(r, p.file.logger, key, 0)
	if err != nil {
		return nil, err
	}

	if len(fds) == 0 {
		return nil, nil
	}

	countFiles := len(fds)

	if p.file.maxFiles > 0 && countFiles > p.file.maxFiles {
		countFiles = p.file.maxFiles
	}

	if err = p.file.checkTotalSize(fds, countFiles); err != nil {
		return nil, err
	}

	images := make([]mrtype.ImageHeader, 0, countFiles)

	for i := 0; i < countFiles; i++ {
		err = func() error { // for defer file.Close()
			if err := p.file.checkFile(fds[i]); err != nil {
				return err
			}

			file, err := fds[i].Open()
			if err != nil {
				return errors.ErrSystemHttpMultipartFormFile.Wrap(err, "key", key)
			}

			defer file.Close()

			contentType := p.file.detectedContentType(fds[i])

			meta, err := p.decode(file, contentType)
			if err != nil {
				return err
			}

			images = append(
				images,
				mrtype.ImageHeader{
					ImageInfo: mrtype.ImageInfo{
						ContentType:  contentType,
						OriginalName: fds[i].Filename,
						Width:        meta.width,
						Height:       meta.height,
						Size:         uint64(fds[i].Size), //nolint:gosec
					},
					Header: fds[i],
				},
			)

			return nil
		}()
		if err != nil {
			return nil, err
		}
	}

	return images, nil
}

func (p *Image) decode(file multipart.File, contentType string) (imageMeta, error) {
	cfg, err := ximage.DecodeImageConfig(file, contentType)
	if err != nil {
		return imageMeta{}, err
	}

	if cfg.Width < 0 || cfg.Height < 0 {
		return imageMeta{}, errors.ErrValidateImageSize
	}

	if p.maxWidth > 0 && uint64(cfg.Width) > p.maxWidth {
		return imageMeta{}, errors.ErrValidateImageWidthMax.New(p.maxWidth)
	}

	if p.maxHeight > 0 && uint64(cfg.Height) > p.maxHeight {
		return imageMeta{}, errors.ErrValidateImageHeightMax.New(p.maxHeight)
	}

	if p.checkBody {
		if err = ximage.CheckImage(file, contentType); err != nil {
			return imageMeta{}, err
		}
	}

	return imageMeta{
		width:  uint64(cfg.Width),
		height: uint64(cfg.Height),
	}, nil
}
