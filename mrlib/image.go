package mrlib

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

// DecodeImageConfig - возвращает информацию об изображении из file, или ошибку, если информацию не удалось извлечь.
func DecodeImageConfig(file io.ReadSeeker, expectedContentType string) (image.Config, error) {
	cfg, err := unsafeDecodeImageConfig(file, expectedContentType)
	if err != nil {
		return image.Config{}, err
	}

	// return offset after DecodeConfig
	_, err = file.Seek(0, 0)

	return cfg, err
}

func unsafeDecodeImageConfig(file io.ReadSeeker, contentType string) (image.Config, error) {
	switch strings.TrimPrefix(contentType, "image/") {
	case "jpg", "jpeg", "jpe":
		return jpeg.DecodeConfig(file)
	case "png":
		return png.DecodeConfig(file)
	case "gif":
		return gif.DecodeConfig(file)
	}

	return image.Config{}, fmt.Errorf("the image content type is not supported: %s", contentType)
}

// CheckImage - возвращает ошибку, если изображение не удалось извлечь из file.
func CheckImage(file io.ReadSeeker, expectedContentType string) error {
	_, err := DecodeImage(file, expectedContentType)

	return err
}

// DecodeImage - возвращает изображение из file, или ошибку, если изображение не удалось извлечь.
func DecodeImage(file io.ReadSeeker, expectedContentType string) (image.Image, error) {
	img, err := unsafeDecodeImage(file, expectedContentType)
	if err != nil {
		return nil, err
	}

	// return offset after Decode
	_, err = file.Seek(0, 0)

	return img, err
}

func unsafeDecodeImage(file io.ReadSeeker, contentType string) (image.Image, error) {
	switch strings.TrimPrefix(contentType, "image/") {
	case "jpg", "jpeg", "jpe":
		return jpeg.Decode(file)
	case "png":
		return png.Decode(file)
	case "gif":
		return gif.Decode(file)
	}

	return nil, fmt.Errorf("the image content type is not supported: %s", contentType)
}
