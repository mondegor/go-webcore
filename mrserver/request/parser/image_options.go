package parser

type (
	// ImageOption - настройка объекта Image.
	ImageOption func(o *imageOptions)

	imageOptions struct {
		image    *Image
		fileOpts []FileOption
	}
)

// WithImageMaxWidth - устанавливает опцию maxWidth для Image (pixels).
func WithImageMaxWidth(value int32) ImageOption {
	return func(o *imageOptions) {
		o.image.maxWidth = value
	}
}

// WithImageMaxHeight - устанавливает опцию maxHeight для Image (pixels).
func WithImageMaxHeight(value int32) ImageOption {
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

// WithImageFileOpts - устанавливает опции File для Image.
func WithImageFileOpts(opts ...FileOption) ImageOption {
	return func(o *imageOptions) {
		o.fileOpts = append(o.fileOpts, opts...)
	}
}
