package parser

import "github.com/mondegor/go-core/util/mime"

type (
	// FileOption - настройка объекта File.
	FileOption func(o *fileOptions)

	fileOptions struct {
		file *File
	}
)

// WithFileMinSize - устанавливает опцию minSize для File (bytes).
func WithFileMinSize(value int64) FileOption {
	return func(o *fileOptions) {
		o.file.minSize = value
	}
}

// WithFileMaxSize - устанавливает опцию maxSize для File (bytes).
func WithFileMaxSize(value int64) FileOption {
	return func(o *fileOptions) {
		o.file.maxSize = value
	}
}

// WithFileMaxFiles - устанавливает опцию maxFiles для File.
func WithFileMaxFiles(value int) FileOption {
	return func(o *fileOptions) {
		o.file.maxFiles = value
	}
}

// WithFileCheckRequestContentType - устанавливает опцию проверки заголовка ContentType в запросе.
func WithFileCheckRequestContentType(value bool) FileOption {
	return func(o *fileOptions) {
		o.file.checkRequestContentType = value
	}
}

// WithFileAllowedMimeTypes - устанавливает опцию с разрешенными типами файлов.
func WithFileAllowedMimeTypes(values []mime.Type) FileOption {
	return func(o *fileOptions) {
		o.file.allowedMimeTypes = mime.NewTypeList(values)
	}
}
