package parser

type (
	// ImageOption - настройка объекта Image.
	ImageOption func(o *imageOptions)

	imageOptions struct {
		image       *Image
		fileOptions []FileOption
	}
)

// WithImageMaxWidth - устанавливает опцию maxWidth для Image (pixels).
func WithImageMaxWidth(value uint64) ImageOption {
	return func(o *imageOptions) {
		o.image.maxWidth = value
	}
}

// WithImageMaxHeight - устанавливает опцию maxHeight для Image (pixels).
func WithImageMaxHeight(value uint64) ImageOption {
	return func(o *imageOptions) {
		o.image.maxHeight = value
	}
}

// WithImageCheckBody - устанавливает опцию checkBody для Image.
func WithImageCheckBody(value bool) ImageOption {
	return func(o *imageOptions) {
		o.image.checkBody = value
	}
}

// WithImageFileOptions - устанавливает опции File для Image.
func WithImageFileOptions(opts ...FileOption) ImageOption {
	return func(o *imageOptions) {
		o.fileOptions = append(o.fileOptions, opts...)
	}
}
