package mrparser

import (
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrserver"
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
		contentType string
		width       int32
		height      int32
	}
)

// Make sure the Image conforms with the mrserver.RequestParserImage interface
var _ mrserver.RequestParserImage = (*Image)(nil)

func NewImage(opts ImageOptions) *Image {
	if len(opts.File.AllowedExts) == 0 {
		opts.File.AllowedExts = []string{".jpeg", ".jpg", ".gif", ".png"}
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
	raw, err := p.file.raw(r, key)

	if err != nil {
		return mrtype.Image{}, err
	}

	meta, err := p.decode(raw.file, raw.contentType)

	if err != nil {
		return mrtype.Image{}, err
	}

	return mrtype.Image{
		ImageInfo: mrtype.ImageInfo{
			ContentType:  meta.contentType,
			OriginalName: raw.hdr.Filename,
			Width:        meta.width,
			Height:       meta.height,
			Size:         raw.hdr.Size,
		},
		Body: raw.file,
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
		contentType: contentType,
		width:       int32(cfg.Width),
		height:      int32(cfg.Height),
	}, nil
}
