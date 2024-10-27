package mrparser

import "github.com/mondegor/go-webcore/mrlib"

type (
	// FileOption - настройка объекта File.
	FileOption func(f *File)
)

// WithFileMinSize - устанавливает опцию minSize для File (bytes).
func WithFileMinSize(value uint64) FileOption {
	return func(f *File) {
		f.minSize = value
	}
}

// WithFileMaxSize - устанавливает опцию maxSize для File (bytes).
func WithFileMaxSize(value uint64) FileOption {
	return func(f *File) {
		f.maxSize = value
	}
}

// WithFileMaxFiles - устанавливает опцию maxFiles для File.
func WithFileMaxFiles(value uint32) FileOption {
	return func(f *File) {
		f.maxFiles = value
	}
}

// WithFileCheckRequestContentType - устанавливает опцию проверки заголовка ContentType в запросе для File.
func WithFileCheckRequestContentType(value bool) FileOption {
	return func(f *File) {
		f.checkRequestContentType = value
	}
}

// WithFileAllowedMimeTypes - устанавливает опцию с разрешенными типами файлов для File.
func WithFileAllowedMimeTypes(values []mrlib.MimeType) FileOption {
	return func(f *File) {
		if len(values) > 0 {
			f.allowedMimeTypes = mrlib.NewMimeTypeList(f.logger, values)
		}
	}
}
