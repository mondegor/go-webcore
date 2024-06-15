package mrtype

import (
	"io"
	"mime/multipart"
	"path"
	"time"
)

type (
	// ImageInfo - comment struct.
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

	// Image - comment struct.
	Image struct {
		ImageInfo
		Body io.ReadCloser
	}

	// ImageContent - comment struct.
	ImageContent struct {
		ImageInfo
		Body []byte
	}

	// ImageHeader - comment struct.
	ImageHeader struct {
		ImageInfo
		Header *multipart.FileHeader
	}
)

// ToFile - comment method.
func (i *ImageInfo) ToFile() FileInfo {
	return FileInfo{
		ContentType:  i.ContentType,
		OriginalName: i.OriginalName,
		Name:         i.Name,
		Path:         i.Path,
		URL:          i.URL,
		Size:         i.Size,
		CreatedAt:    TimePointerCopy(i.CreatedAt),
		UpdatedAt:    TimePointerCopy(i.UpdatedAt),
	}
}

// Original - comment method.
func (i *ImageInfo) Original() string {
	if i.OriginalName != "" {
		return i.OriginalName
	}

	if i.Name != "" {
		return i.Name
	}

	return path.Base(i.Path)
}

// ToFile - comment method.
func (i *Image) ToFile() File {
	return File{
		FileInfo: i.ImageInfo.ToFile(),
		Body:     i.Body,
	}
}

// ToFile - comment method.
func (i *ImageContent) ToFile() FileContent {
	return FileContent{
		FileInfo: i.ImageInfo.ToFile(),
		Body:     i.Body,
	}
}

// ToFile - comment method.
func (i *ImageHeader) ToFile() FileHeader {
	return FileHeader{
		FileInfo: i.ImageInfo.ToFile(),
		Header:   i.Header,
	}
}
