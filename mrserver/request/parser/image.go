package parser

import (
	"bytes"
	"math"
	"mime/multipart"
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrmodel/media"
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
		maxWidth  int32 // pixels
		maxHeight int32 // pixels
		checkBody bool
		file      *File
	}

	// imageMeta - метаданные изображения (ширина и высота).
	imageMeta struct {
		width  int32
		height int32
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
		fileOpts: []FileOption{
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

	o.image.file = NewFile(logger, o.fileOpts...)

	return o.image
}

// FormImage - возвращает информацию об изображении со ссылкой для чтения файла изображения из MultipartForm.
// WARNING: you don't forget to call result.Body.Close().
func (p *Image) FormImage(r *http.Request, key string) (media.Image, error) {
	hdr, err := p.file.formFile(r, p.file.logger, key)
	if err != nil {
		return media.Image{}, err
	}

	if err = p.file.checkFile(hdr); err != nil {
		return media.Image{}, err
	}

	file, err := hdr.Open()
	if err != nil {
		return media.Image{}, errors.ErrSystemHttpMultipartFormFile.Wrap(err, "key", key)
	}

	contentType := p.file.detectedContentType(hdr)

	meta, err := p.decode(file, contentType)
	if err != nil {
		if err := file.Close(); err != nil {
			p.file.logger.Error(r.Context(), "mrparser.FormImage: error when closing file image", "error", err)
		}

		return media.Image{}, err
	}

	return media.Image{
		ImageInfo: media.ImageInfo{
			ContentType:  contentType,
			OriginalName: hdr.Filename,
			Width:        meta.width,
			Height:       meta.height,
			Size:         hdr.Size,
		},
		Body: file,
	}, nil
}

// FormImageContent - возвращает информацию об изображении и сам файл изображения из MultipartForm.
// WARNING: only for short files.
func (p *Image) FormImageContent(r *http.Request, key string) (media.ImageContent, error) {
	file, err := p.FormImage(r, key)
	if err != nil {
		return media.ImageContent{}, err
	}

	defer func() {
		_ = file.Body.Close()
	}()

	var buf bytes.Buffer

	if _, err = buf.ReadFrom(file.Body); err != nil {
		return media.ImageContent{}, errors.WrapInternalError(err, "reading file.Body failed")
	}

	return media.ImageContent{
		ImageInfo: file.ImageInfo,
		Body:      buf.Bytes(),
	}, nil
}

// FormImages - возвращает массив заголовков на файлы изображений из MultipartForm.
func (p *Image) FormImages(r *http.Request, key string) ([]media.ImageHeader, error) {
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

	images := make([]media.ImageHeader, 0, countFiles)

	for i := 0; i < countFiles; i++ {
		err = func() error { // for defer file.Close()
			if err := p.file.checkFile(fds[i]); err != nil {
				return err
			}

			file, err := fds[i].Open()
			if err != nil {
				return errors.ErrSystemHttpMultipartFormFile.Wrap(err, "key", key)
			}

			defer func() {
				_ = file.Close()
			}()

			contentType := p.file.detectedContentType(fds[i])

			meta, err := p.decode(file, contentType)
			if err != nil {
				return err
			}

			images = append(
				images,
				media.ImageHeader{
					ImageInfo: media.ImageInfo{
						ContentType:  contentType,
						OriginalName: fds[i].Filename,
						Width:        meta.width,
						Height:       meta.height,
						Size:         fds[i].Size,
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

	if cfg.Width < 0 || cfg.Width > math.MinInt32 || cfg.Height < 0 || cfg.Height > math.MinInt32 {
		return imageMeta{}, errors.ErrValidateImageSize
	}

	if p.maxWidth > 0 && cfg.Width > int(p.maxWidth) {
		return imageMeta{}, errors.ErrValidateImageWidthMax.New(p.maxWidth)
	}

	if p.maxHeight > 0 && cfg.Height > int(p.maxHeight) {
		return imageMeta{}, errors.ErrValidateImageHeightMax.New(p.maxHeight)
	}

	if p.checkBody {
		if err = ximage.CheckImage(file, contentType); err != nil {
			return imageMeta{}, err
		}
	}

	return imageMeta{
		width:  int32(cfg.Width),
		height: int32(cfg.Height),
	}, nil
}
