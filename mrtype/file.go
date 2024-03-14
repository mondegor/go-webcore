package mrtype

import (
	"io"
	"mime/multipart"
	"path"
	"time"
)

type (
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

	File struct {
		FileInfo
		Body io.ReadCloser
	}

	FileContent struct {
		FileInfo
		Body []byte
	}

	FileHeader struct {
		FileInfo
		Header *multipart.FileHeader
	}
)

func (f *FileInfo) Original() string {
	if f.OriginalName != "" {
		return f.OriginalName
	}

	if f.Name != "" {
		return f.Name
	}

	return path.Base(f.Path)
}
