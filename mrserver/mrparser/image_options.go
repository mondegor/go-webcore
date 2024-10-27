package mrparser

type (
	// ImageOption - настройка объекта Image.
	ImageOption func(i *Image)
)

// WithImageMaxWidth - устанавливает опцию maxWidth для Image (pixels).
func WithImageMaxWidth(value uint64) ImageOption {
	return func(i *Image) {
		i.maxWidth = value
	}
}

// WithImageMaxHeight - устанавливает опцию maxHeight для Image (pixels).
func WithImageMaxHeight(value uint64) ImageOption {
	return func(i *Image) {
		i.maxHeight = value
	}
}

// WithImageCheckBody - устанавливает опцию maxHeight для Image.
func WithImageCheckBody(value bool) ImageOption {
	return func(i *Image) {
		i.checkBody = value
	}
}

// WithImageFileOptions - устанавливает опции File для Image.
func WithImageFileOptions(opts ...FileOption) ImageOption {
	return func(i *Image) {
		i.file.applyOptions(opts)
	}
}
