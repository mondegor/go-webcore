package mrparser

import (
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
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
				[]mrlib.MimeType{
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
	hdr, err := mrreq.FormFile(r, key)
	if err != nil {
		return mrtype.Image{}, err
	}

	if err = p.file.checkFile(hdr); err != nil {
		return mrtype.Image{}, err
	}

	file, err := hdr.Open()
	if err != nil {
		return mrtype.Image{}, mrcore.ErrHttpMultipartFormFile.Wrap(err, key)
	}

	contentType := p.file.detectedContentType(hdr)

	meta, err := p.decode(file, contentType)
	if err != nil {
		if err := file.Close(); err != nil {
			mrlog.Ctx(r.Context()).Error().Err(err).Msg("mrparser.FormImage: error when closing file image")
		}

		return mrtype.Image{}, err
	}

	return mrtype.Image{
		ImageInfo: mrtype.ImageInfo{
			ContentType:  contentType,
			OriginalName: hdr.Filename,
			Width:        meta.width,
			Height:       meta.height,
			Size:         uint64(hdr.Size),
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
		return mrtype.ImageContent{}, mrcore.ErrInternal.Wrap(err)
	}

	return mrtype.ImageContent{
		ImageInfo: file.ImageInfo,
		Body:      buf.Bytes(),
	}, nil
}

// FormImages - возвращает массив заголовков на файлы изображений из MultipartForm.
func (p *Image) FormImages(r *http.Request, key string) ([]mrtype.ImageHeader, error) {
	fds, err := mrreq.FormFiles(r, key, 0)
	if err != nil {
		return nil, err
	}

	if len(fds) == 0 {
		return nil, nil
	}

	countFiles := p.file.allowedFiles(uint32(len(fds)))

	if err = p.file.checkTotalSize(fds, countFiles); err != nil {
		return nil, err
	}

	images := make([]mrtype.ImageHeader, 0, countFiles)

	for i := uint32(0); i < countFiles; i++ {
		err = func() error { // for defer file.Close()
			if err := p.file.checkFile(fds[i]); err != nil {
				return err
			}

			file, err := fds[i].Open()
			if err != nil {
				return mrcore.ErrHttpMultipartFormFile.Wrap(err, key)
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
						Size:         uint64(fds[i].Size),
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
	cfg, err := mrlib.DecodeImageConfig(file, contentType)
	if err != nil {
		return imageMeta{}, err
	}

	if p.maxWidth > 0 && uint64(cfg.Width) > p.maxWidth {
		return imageMeta{}, ErrHttpRequestImageWidthMax.New(p.maxWidth)
	}

	if p.maxHeight > 0 && uint64(cfg.Height) > p.maxHeight {
		return imageMeta{}, ErrHttpRequestImageHeightMax.New(p.maxHeight)
	}

	if p.checkBody {
		if err = mrlib.CheckImage(file, contentType); err != nil {
			return imageMeta{}, err
		}
	}

	return imageMeta{
		width:  uint64(cfg.Width),
		height: uint64(cfg.Height),
	}, nil
}
