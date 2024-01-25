package mrparser

import (
	. "github.com/mondegor/go-sysmess/mrerr"
)

var (
	FactoryErrHttpRequestFileSizeMin = NewFactory(
		"errHttpRequestFileSizeMin", ErrorKindUser, "invalid file size, min size = {{ .value }}b")

	FactoryErrHttpRequestFileSizeMax = NewFactory(
		"errHttpRequestFileSizeMax", ErrorKindUser, "invalid file size, max size = {{ .value }}b")

	FactoryErrHttpRequestFileExtension = NewFactory(
		"errHttpRequestFileExtension", ErrorKindUser, "invalid file extension: {{ .value }}")

	FactoryErrHttpRequestFileContentType = NewFactory(
		"errHttpRequestFileContentType", ErrorKindUser, "the content type '{{ .value }}' does not match the detected type")

	FactoryErrHttpRequestFileUnsupportedType = NewFactory(
		"errHttpRequestFileUnsupportedType", ErrorKindUser, "unsupported file type '{{ .value }}'")

	FactoryErrHttpRequestImageWidthMax = NewFactory(
		"errHttpRequestImageWidthMax", ErrorKindUser, "invalid image width, max size = {{ .value }}px")

	FactoryErrHttpRequestImageHeightMax = NewFactory(
		"errHttpRequestImageHeightMax", ErrorKindUser, "invalid image height, max size = {{ .value }}px")
)

func WrapFileError(err error, name string) error {
	if FactoryErrHttpRequestFileSizeMin.Is(err) {
		return NewCustomError(name, err)
	}

	if FactoryErrHttpRequestFileSizeMax.Is(err) {
		return NewCustomError(name, err)
	}

	if FactoryErrHttpRequestFileExtension.Is(err) {
		return NewCustomError(name, err)
	}

	if FactoryErrHttpRequestFileContentType.Is(err) {
		return NewCustomError(name, err)
	}

	if FactoryErrHttpRequestFileUnsupportedType.Is(err) {
		return NewCustomError(name, err)
	}

	return err
}

func WrapImageError(err error, name string) error {
	if FactoryErrHttpRequestImageWidthMax.Is(err) {
		return NewCustomError(name, err)
	}

	if FactoryErrHttpRequestImageHeightMax.Is(err) {
		return NewCustomError(name, err)
	}

	return WrapFileError(err, name)
}
