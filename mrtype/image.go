package mrtype

import (
	"io"
	"time"
)

type (
	ImageInfo struct {
		ContentType  string     `json:"contentType,omitempty"`
		OriginalName string     `json:"originalName,omitempty"`
		Name         string     `json:"name,omitempty"`
		Path         string     `json:"-"`
		URL          string     `json:"url,omitempty"`
		Width        int32      `json:"width,omitempty"`
		Height       int32      `json:"height,omitempty"`
		Size         int64      `json:"size,omitempty"`
		CreatedAt    *time.Time `json:"created,omitempty"`
		ModifiedAt   *time.Time `json:"lastModified,omitempty"`
	}

	Image struct {
		ImageInfo
		Body io.ReadCloser
	}

	ImageContent struct {
		ImageInfo
		Body []byte
	}
)

func (i *ImageInfo) ToFile() FileInfo {
	return FileInfo{
		ContentType:  i.ContentType,
		OriginalName: i.OriginalName,
		Name:         i.Name,
		Path:         i.Path,
		URL:          i.URL,
		Size:         i.Size,
		CreatedAt:    TimePointerCopy(i.CreatedAt),
		ModifiedAt:   TimePointerCopy(i.ModifiedAt),
	}
}

func (i *Image) ToFile() File {
	return File{
		FileInfo: i.ImageInfo.ToFile(),
		Body:     i.Body,
	}
}

func (i *ImageContent) ToFile() FileContent {
	return FileContent{
		FileInfo: i.ImageInfo.ToFile(),
		Body:     i.Body,
	}
}
