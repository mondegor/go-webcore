package mrtype

import (
	"io"
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
		CreatedAt    *time.Time `json:"created,omitempty"`
		ModifiedAt   *time.Time `json:"lastModified,omitempty"`
	}

	File struct {
		FileInfo
		Body io.ReadCloser
	}

	FileContent struct {
		FileInfo
		Body []byte
	}
)
