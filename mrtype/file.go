package mrtype

import (
	"io"
	"mime/multipart"
	"path"
	"time"
)

type (
	// FileInfo - мета-информация о файле.
	FileInfo struct {
		ContentType  string     `json:"contentType,omitempty"`
		OriginalName string     `json:"originalName,omitempty"`
		Name         string     `json:"name,omitempty"`
		Path         string     `json:"-"`
		URL          string     `json:"url,omitempty"`
		Size         int64      `json:"size,omitempty"`
		CreatedAt    *time.Time `json:"createdAt,omitempty"`
		UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
	}

	// File - мета-информация файла вместе с источником файла.
	File struct {
		FileInfo
		Body io.ReadCloser
	}

	// FileContent - файл с мета-информацией.
	FileContent struct {
		FileInfo
		Body []byte
	}

	// FileHeader - мета-информация файла вместе с источником файла (multipart/form-data).
	FileHeader struct {
		FileInfo
		Header *multipart.FileHeader
	}
)

// Original - возвращает оригинальное имя файла (как оно было названо в первоисточнике).
func (f *FileInfo) Original() string {
	if f.OriginalName != "" {
		return f.OriginalName
	}

	if f.Name != "" {
		return f.Name
	}

	return path.Base(f.Path)
}
