package mrparser

import (
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Image struct {
		file      *File
		maxWidth  int32
		maxHeight int32
		checkBody bool
	}

	ImageOptions struct {
		File      FileOptions
		MaxWidth  int32
		MaxHeight int32
		CheckBody bool
	}

	imageMeta struct {
		width  int32
		height int32
	}
)

// Make sure the Image conforms with the mrserver.RequestParserImage interface
var _ mrserver.RequestParserImage = (*Image)(nil)

func NewImage(opts ImageOptions) *Image {
	if len(opts.File.AllowedExts) == 0 {
		opts.File.AllowedExts = []string{".gif", ".jpeg", ".jpg", ".png"}
	}

	return &Image{
		file:      NewFile(opts.File),
		maxWidth:  opts.MaxWidth,
		maxHeight: opts.MaxHeight,
		checkBody: opts.CheckBody,
	}
}

// FormImage - WARNING you don't forget to call result.Body.Close()
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
		return mrtype.Image{}, mrcore.FactoryErrHttpMultipartFormFile.Wrap(err, key)
	}

	contentType := p.file.detectedContentType(hdr)
	meta, err := p.decode(file, contentType)

	if err != nil {
		file.Close()
		return mrtype.Image{}, err
	}

	return mrtype.Image{
		ImageInfo: mrtype.ImageInfo{
			ContentType:  contentType,
			OriginalName: hdr.Filename,
			Width:        meta.width,
			Height:       meta.height,
			Size:         hdr.Size,
		},
		Body: file,
	}, nil
}

// FormImageContent - only for short files
func (p *Image) FormImageContent(r *http.Request, key string) (mrtype.ImageContent, error) {
	file, err := p.FormImage(r, key)

	if err != nil {
		return mrtype.ImageContent{}, nil
	}

	defer file.Body.Close()

	buf := bytes.Buffer{}

	if _, err = buf.ReadFrom(file.Body); err != nil {
		return mrtype.ImageContent{}, mrcore.FactoryErrInternal.Wrap(err)
	}

	return mrtype.ImageContent{
		ImageInfo: file.ImageInfo,
		Body:      buf.Bytes(),
	}, nil
}

func (p *Image) FormImages(r *http.Request, key string) ([]mrtype.ImageHeader, error) {
	fds, err := mrreq.FormFiles(r, key, 0)

	if len(fds) == 0 {
		return nil, nil
	}

	countFiles := p.file.allowedFiles(len(fds))

	if err = p.file.checkTotalSize(fds, countFiles); err != nil {
		return nil, err
	}

	images := make([]mrtype.ImageHeader, countFiles)

	for i := 0; i < countFiles; i++ {
		err = func() error { // for defer file.Close()
			if err := p.file.checkFile(fds[i]); err != nil {
				return err
			}

			file, err := fds[i].Open()

			if err != nil {
				return mrcore.FactoryErrHttpMultipartFormFile.Wrap(err, key)
			}

			defer file.Close()

			contentType := p.file.detectedContentType(fds[i])
			meta, err := p.decode(file, contentType)

			if err != nil {
				return err
			}

			images[i] = mrtype.ImageHeader{
				ImageInfo: mrtype.ImageInfo{
					ContentType:  contentType,
					OriginalName: fds[i].Filename,
					Width:        meta.width,
					Height:       meta.height,
					Size:         fds[i].Size,
				},
				Header: fds[i],
			}

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

	if p.maxWidth > 0 && int32(cfg.Width) > p.maxWidth {
		return imageMeta{}, FactoryErrHttpRequestImageWidthMax.New(p.maxWidth)
	}

	if p.maxHeight > 0 && int32(cfg.Height) > p.maxHeight {
		return imageMeta{}, FactoryErrHttpRequestImageHeightMax.New(p.maxHeight)
	}

	if p.checkBody {
		if err = mrlib.CheckImage(file, contentType); err != nil {
			return imageMeta{}, err
		}
	}

	return imageMeta{
		width:  int32(cfg.Width),
		height: int32(cfg.Height),
	}, nil
}
