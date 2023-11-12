package mrtype

import (
	"io"
	"time"
)

type (
	FileInfo struct {
		ContentType  string	`json:"contentType,omitempty"`
		OriginalName string	`json:"originalName,omitempty"`
		Name		 string	`json:"name"`
		LastModified time.Time `json:"-"`
		Size		 int64	 `json:"size"`
	}

	File struct {
		FileInfo
		Path string
		Body io.ReadCloser
	}
)
