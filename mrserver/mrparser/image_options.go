package mrparser

type (
	// ImageOption - настройка объекта Image.
	ImageOption func(im *Image)
)

// WithImageMaxWidth - устанавливает опцию maxWidth для Image (pixels).
func WithImageMaxWidth(value uint64) ImageOption {
	return func(im *Image) {
		im.maxWidth = value
	}
}

// WithImageMaxHeight - устанавливает опцию maxHeight для Image (pixels).
func WithImageMaxHeight(value uint64) ImageOption {
	return func(im *Image) {
		im.maxHeight = value
	}
}

// WithImageCheckBody - устанавливает опцию maxHeight для Image.
func WithImageCheckBody(value bool) ImageOption {
	return func(im *Image) {
		im.checkBody = value
	}
}

// WithImageFileOptions - устанавливает опции File для Image.
func WithImageFileOptions(opts ...FileOption) ImageOption {
	return func(im *Image) {
		im.file.applyOptions(opts)
	}
}
