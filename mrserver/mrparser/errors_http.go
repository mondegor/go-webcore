package mrparser

import "github.com/mondegor/go-sysmess/mrerr"

var (
	// ErrHttpRequestFileSizeMin - invalid file size - min.
	ErrHttpRequestFileSizeMin = mrerr.NewProto(
		"errHttpRequestFileSizeMin", mrerr.ErrorKindUser, "invalid file size, min size = {{ .value }}b")

	// ErrHttpRequestFileSizeMax - invalid file size - max.
	ErrHttpRequestFileSizeMax = mrerr.NewProto(
		"errHttpRequestFileSizeMax", mrerr.ErrorKindUser, "invalid file size, max size = {{ .value }}b")

	// ErrHttpRequestFileExtension - invalid file extension.
	ErrHttpRequestFileExtension = mrerr.NewProto(
		"errHttpRequestFileExtension", mrerr.ErrorKindUser, "invalid file extension: {{ .value }}")

	// ErrHttpRequestFileTotalSizeMax - invalid file total size - max.
	ErrHttpRequestFileTotalSizeMax = mrerr.NewProto(
		"errHttpRequestFileTotalSizeMax", mrerr.ErrorKindUser, "invalid file total size, max total size = {{ .value }}b")

	// ErrHttpRequestFileContentType - the content type does not match the detected type.
	ErrHttpRequestFileContentType = mrerr.NewProto(
		"errHttpRequestFileContentType", mrerr.ErrorKindUser, "the content type '{{ .value }}' does not match the detected type")

	// ErrHttpRequestFileUnsupportedType - unsupported file type.
	ErrHttpRequestFileUnsupportedType = mrerr.NewProto(
		"errHttpRequestFileUnsupportedType", mrerr.ErrorKindUser, "unsupported file type '{{ .value }}'")

	// ErrHttpRequestImageWidthMax - invalid image width - max.
	ErrHttpRequestImageWidthMax = mrerr.NewProto(
		"errHttpRequestImageWidthMax", mrerr.ErrorKindUser, "invalid image width, max size = {{ .value }}px")

	// ErrHttpRequestImageHeightMax - invalid image height - max.
	ErrHttpRequestImageHeightMax = mrerr.NewProto(
		"errHttpRequestImageHeightMax", mrerr.ErrorKindUser, "invalid image height, max size = {{ .value }}px")
)

// WrapFileError - comment func.
func WrapFileError(err error, name string) error {
	if ErrHttpRequestFileSizeMin.Is(err) {
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestFileSizeMax.Is(err) {
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestFileExtension.Is(err) {
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestFileContentType.Is(err) {
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestFileUnsupportedType.Is(err) {
		return mrerr.NewCustomError(name, err)
	}

	return err
}

// WrapImageError - comment func.
func WrapImageError(err error, name string) error {
	if ErrHttpRequestImageWidthMax.Is(err) {
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestImageHeightMax.Is(err) {
		return mrerr.NewCustomError(name, err)
	}

	return WrapFileError(err, name)
}
