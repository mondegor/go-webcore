package mrtype

import (
	"io"
	"mime/multipart"
	"path"
	"time"
)

type (
	// ImageInfo - мета-информация об изображении.
	ImageInfo struct {
		ContentType  string     `json:"contentType,omitempty"`
		OriginalName string     `json:"originalName,omitempty"`
		Name         string     `json:"name,omitempty"`
		Path         string     `json:"-"`
		URL          string     `json:"url,omitempty"`
		Width        int32      `json:"width,omitempty"`
		Height       int32      `json:"height,omitempty"`
		Size         int64      `json:"size,omitempty"`
		CreatedAt    *time.Time `json:"createdAt,omitempty"`
		UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
	}

	// Image - мета-информация изображения вместе с источником изображения.
	Image struct {
		ImageInfo
		Body io.ReadCloser
	}

	// ImageContent - изображение с мета-информацией.
	ImageContent struct {
		ImageInfo
		Body []byte
	}

	// ImageHeader - мета-информация изображения вместе с источником изображения (multipart/form-data).
	ImageHeader struct {
		ImageInfo
		Header *multipart.FileHeader
	}
)

// Original - возвращает оригинальное имя изображение (как оно было названо в первоисточнике).
func (i *ImageInfo) Original() string {
	if i.OriginalName != "" {
		return i.OriginalName
	}

	if i.Name != "" {
		return i.Name
	}

	return path.Base(i.Path)
}

// ToFileInfo - возвращает мета-информацию изображения преобразованное в файловую структуру
// (с потерей дополнительной информации об изображении).
func (i *ImageInfo) ToFileInfo() FileInfo {
	return FileInfo{
		ContentType:  i.ContentType,
		OriginalName: i.OriginalName,
		Name:         i.Name,
		Path:         i.Path,
		URL:          i.URL,
		Size:         i.Size,
		CreatedAt:    CopyTimePointer(i.CreatedAt),
		UpdatedAt:    CopyTimePointer(i.UpdatedAt),
	}
}

// ToFile - возвращает изображение преобразованное в файловую структуру
// (с потерей дополнительной информации об изображении).
func (i *Image) ToFile() File {
	return File{
		FileInfo: i.ImageInfo.ToFileInfo(),
		Body:     i.Body,
	}
}

// ToFileContent - возвращает изображение преобразованное в файловую структуру
// (с потерей дополнительной информации об изображении).
func (i *ImageContent) ToFileContent() FileContent {
	return FileContent{
		FileInfo: i.ImageInfo.ToFileInfo(),
		Body:     i.Body,
	}
}

// ToFileHeader - возвращает изображение преобразованное в файловую структуру
// (с потерей дополнительной информации об изображении).
func (i *ImageHeader) ToFileHeader() FileHeader {
	return FileHeader{
		FileInfo: i.ImageInfo.ToFileInfo(),
		Header:   i.Header,
	}
}
