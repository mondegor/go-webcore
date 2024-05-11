package mrparser

import (
	e "github.com/mondegor/go-sysmess/mrerr"
)

var (
	FactoryErrHTTPRequestFileSizeMin = e.NewFactory(
		"errHttpRequestFileSizeMin", e.ErrorTypeUser, "invalid file size, min size = {{ .value }}b")

	FactoryErrHTTPRequestFileSizeMax = e.NewFactory(
		"errHttpRequestFileSizeMax", e.ErrorTypeUser, "invalid file size, max size = {{ .value }}b")

	FactoryErrHTTPRequestFileExtension = e.NewFactory(
		"errHttpRequestFileExtension", e.ErrorTypeUser, "invalid file extension: {{ .value }}")

	FactoryErrHTTPRequestFileTotalSizeMax = e.NewFactory(
		"errHttpRequestFileTotalSizeMax", e.ErrorTypeUser, "invalid file total size, max total size = {{ .value }}b")

	FactoryErrHTTPRequestFileContentType = e.NewFactory(
		"errHttpRequestFileContentType", e.ErrorTypeUser, "the content type '{{ .value }}' does not match the detected type")

	FactoryErrHTTPRequestFileUnsupportedType = e.NewFactory(
		"errHttpRequestFileUnsupportedType", e.ErrorTypeUser, "unsupported file type '{{ .value }}'")

	FactoryErrHTTPRequestImageWidthMax = e.NewFactory(
		"errHttpRequestImageWidthMax", e.ErrorTypeUser, "invalid image width, max size = {{ .value }}px")

	FactoryErrHTTPRequestImageHeightMax = e.NewFactory(
		"errHttpRequestImageHeightMax", e.ErrorTypeUser, "invalid image height, max size = {{ .value }}px")
)

func WrapFileError(err error, name string) error {
	if FactoryErrHTTPRequestFileSizeMin.Is(err) {
		return e.NewCustomError(name, err)
	}

	if FactoryErrHTTPRequestFileSizeMax.Is(err) {
		return e.NewCustomError(name, err)
	}

	if FactoryErrHTTPRequestFileExtension.Is(err) {
		return e.NewCustomError(name, err)
	}

	if FactoryErrHTTPRequestFileContentType.Is(err) {
		return e.NewCustomError(name, err)
	}

	if FactoryErrHTTPRequestFileUnsupportedType.Is(err) {
		return e.NewCustomError(name, err)
	}

	return err
}

func WrapImageError(err error, name string) error {
	if FactoryErrHTTPRequestImageWidthMax.Is(err) {
		return e.NewCustomError(name, err)
	}

	if FactoryErrHTTPRequestImageHeightMax.Is(err) {
		return e.NewCustomError(name, err)
	}

	return WrapFileError(err, name)
}
