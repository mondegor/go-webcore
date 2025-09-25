package mrparser

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrHttpRequestFileSize - invalid file size.
	ErrHttpRequestFileSize = mrerr.NewKindUser("RequestFileSize", "invalid file size")

	// ErrHttpRequestFileSizeMin - invalid file size - min.
	ErrHttpRequestFileSizeMin = mrerr.NewKindUser("RequestFileSizeMin", "invalid file size, min size = {Value}b")

	// ErrHttpRequestFileSizeMax - invalid file size - max.
	ErrHttpRequestFileSizeMax = mrerr.NewKindUser("RequestFileSizeMax", "invalid file size, max size = {Value}b")

	// ErrHttpRequestFileExtension - invalid file extension.
	ErrHttpRequestFileExtension = mrerr.NewKindUser("RequestFileExtension", "invalid file extension: {Value}")

	// ErrHttpRequestFileTotalSizeMax - invalid file total size - max.
	ErrHttpRequestFileTotalSizeMax = mrerr.NewKindUser("RequestFileTotalSizeMax", "invalid file total size, max total size = {Value}b")

	// ErrHttpRequestFileContentType - the content type does not match the detected type.
	ErrHttpRequestFileContentType = mrerr.NewKindUser("RequestFileContentType", "the content type '{Value}' does not match the detected type")

	// ErrHttpRequestFileUnsupportedType - unsupported file type.
	ErrHttpRequestFileUnsupportedType = mrerr.NewKindUser("RequestFileUnsupportedType", "unsupported file type '{Value}'")

	// ErrHttpRequestImageSize - invalid image size (width, height).
	ErrHttpRequestImageSize = mrerr.NewKindUser("ErrHttpRequestImageSize", "invalid image size (width, height)")

	// ErrHttpRequestImageWidthMax - invalid image width - max.
	ErrHttpRequestImageWidthMax = mrerr.NewKindUser("RequestImageWidthMax", "invalid image width, max size = {Value}px")

	// ErrHttpRequestImageHeightMax - invalid image height - max.
	ErrHttpRequestImageHeightMax = mrerr.NewKindUser("RequestImageHeightMax", "invalid image height, max size = {Value}px")
)

// WrapFileError - оборачивает ошибки связанные с парсингом файла.
func WrapFileError(err error, name string) error {
	if ErrHttpRequestFileSizeMin.Is(err) { // вложенные ошибки не учитываются
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestFileSizeMax.Is(err) { // вложенные ошибки не учитываются
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestFileExtension.Is(err) { // вложенные ошибки не учитываются
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestFileContentType.Is(err) { // вложенные ошибки не учитываются
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestFileUnsupportedType.Is(err) { // вложенные ошибки не учитываются
		return mrerr.NewCustomError(name, err)
	}

	return err
}

// WrapImageError - оборачивает ошибки связанные с парсингом изображения.
func WrapImageError(err error, name string) error {
	if ErrHttpRequestImageWidthMax.Is(err) { // вложенные ошибки не учитываются
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestImageHeightMax.Is(err) { // вложенные ошибки не учитываются
		return mrerr.NewCustomError(name, err)
	}

	return WrapFileError(err, name)
}
