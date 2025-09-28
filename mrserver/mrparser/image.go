package mrparser

import (
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/mondegor/go-sysmess/mrdto"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlib/extfile"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
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
	im := &Image{
		maxWidth:  defaultMaxWidth,
		maxHeight: defaultMaxHeight,
		checkBody: defaultCheckBody,
		file: NewFile( // by default
			logger,
			WithFileAllowedMimeTypes(
				[]extfile.MimeType{
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
		),
	}

	for _, opt := range opts {
		opt(im)
	}

	return im
}

// FormImage - возвращает информацию об изображении со ссылкой для чтения файла изображения из MultipartForm.
// WARNING: you don't forget to call result.Body.Close().
func (p *Image) FormImage(r *http.Request, key string) (mrtype.Image, error) {
	hdr, err := mrreq.FormFile(r, p.file.logger, key)
	if err != nil {
		return mrtype.Image{}, err
	}

	if err = p.file.checkFile(hdr); err != nil {
		return mrtype.Image{}, err
	}

	file, err := hdr.Open()
	if err != nil {
		return mrtype.Image{}, mr.ErrHttpMultipartFormFile.Wrap(err, key)
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
		ImageInfo: mrdto.ImageInfo{
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
		return mrtype.ImageContent{}, mr.ErrInternal.Wrap(err)
	}

	return mrtype.ImageContent{
		ImageInfo: file.ImageInfo,
		Body:      buf.Bytes(),
	}, nil
}

// FormImages - возвращает массив заголовков на файлы изображений из MultipartForm.
func (p *Image) FormImages(r *http.Request, key string) ([]mrtype.ImageHeader, error) {
	fds, err := mrreq.FormFiles(r, p.file.logger, key, 0)
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
				return mr.ErrHttpMultipartFormFile.Wrap(err, key)
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
					ImageInfo: mrdto.ImageInfo{
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
	cfg, err := extfile.DecodeImageConfig(file, contentType)
	if err != nil {
		return imageMeta{}, err
	}

	if cfg.Width < 0 || cfg.Height < 0 {
		return imageMeta{}, mr.ErrValidateImageSize.New()
	}

	if p.maxWidth > 0 && uint64(cfg.Width) > p.maxWidth {
		return imageMeta{}, mr.ErrValidateImageWidthMax.New(p.maxWidth)
	}

	if p.maxHeight > 0 && uint64(cfg.Height) > p.maxHeight {
		return imageMeta{}, mr.ErrValidateImageHeightMax.New(p.maxHeight)
	}

	if p.checkBody {
		if err = extfile.CheckImage(file, contentType); err != nil {
			return imageMeta{}, err
		}
	}

	return imageMeta{
		width:  uint64(cfg.Width),
		height: uint64(cfg.Height),
	}, nil
}
