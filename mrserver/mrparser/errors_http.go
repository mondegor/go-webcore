package mrparser

import "github.com/mondegor/go-sysmess/mrerr"

var (
	// ErrHttpRequestFileSizeMin - comment error.
	ErrHttpRequestFileSizeMin = mrerr.NewProto(
		"errHttpRequestFileSizeMin", mrerr.ErrorKindUser, "invalid file size, min size = {{ .value }}b")

	// ErrHttpRequestFileSizeMax - comment error.
	ErrHttpRequestFileSizeMax = mrerr.NewProto(
		"errHttpRequestFileSizeMax", mrerr.ErrorKindUser, "invalid file size, max size = {{ .value }}b")

	// ErrHttpRequestFileExtension - comment error.
	ErrHttpRequestFileExtension = mrerr.NewProto(
		"errHttpRequestFileExtension", mrerr.ErrorKindUser, "invalid file extension: {{ .value }}")

	// ErrHttpRequestFileTotalSizeMax - comment error.
	ErrHttpRequestFileTotalSizeMax = mrerr.NewProto(
		"errHttpRequestFileTotalSizeMax", mrerr.ErrorKindUser, "invalid file total size, max total size = {{ .value }}b")

	// ErrHttpRequestFileContentType - comment error.
	ErrHttpRequestFileContentType = mrerr.NewProto(
		"errHttpRequestFileContentType", mrerr.ErrorKindUser, "the content type '{{ .value }}' does not match the detected type")

	// ErrHttpRequestFileUnsupportedType - comment error.
	ErrHttpRequestFileUnsupportedType = mrerr.NewProto(
		"errHttpRequestFileUnsupportedType", mrerr.ErrorKindUser, "unsupported file type '{{ .value }}'")

	// ErrHttpRequestImageWidthMax - comment error.
	ErrHttpRequestImageWidthMax = mrerr.NewProto(
		"errHttpRequestImageWidthMax", mrerr.ErrorKindUser, "invalid image width, max size = {{ .value }}px")

	// ErrHttpRequestImageHeightMax - comment error.
	ErrHttpRequestImageHeightMax = mrerr.NewProto(
		"errHttpRequestImageHeightMax", mrerr.ErrorKindUser, "invalid image height, max size = {{ .value }}px")
)

// WrapFileError  - comment func.
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

// WrapImageError  - comment func.
func WrapImageError(err error, name string) error {
	if ErrHttpRequestImageWidthMax.Is(err) {
		return mrerr.NewCustomError(name, err)
	}

	if ErrHttpRequestImageHeightMax.Is(err) {
		return mrerr.NewCustomError(name, err)
	}

	return WrapFileError(err, name)
}
