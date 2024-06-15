package mrlib

import (
	"errors"
	"fmt"
	"path"
)

type (
	// MimeTypeList - comment struct.
	MimeTypeList struct {
		extMap         map[string]string
		contentTypeMap map[string]string
	}

	// MimeType - comment struct.
	MimeType struct {
		Extension   string `yaml:"ext"`
		ContentType string `yaml:"type"`
	}
)

// NewMimeTypeList - создаёт объект MimeTypeList.
func NewMimeTypeList(items []MimeType) *MimeTypeList {
	extMap := make(map[string]string, len(items))
	mimeMap := make(map[string]string, len(items))

	for _, item := range items {
		extMap[item.Extension] = item.ContentType

		// т.к. у одного типа может быть несколько расширений,
		// то в индекс попадает только первый зарегистрированный
		if _, ok := mimeMap[item.ContentType]; !ok {
			mimeMap[item.ContentType] = item.Extension
		}
	}

	return &MimeTypeList{
		extMap:         extMap,
		contentTypeMap: mimeMap,
	}
}

// NewListByExts  - comment func.
func (mt *MimeTypeList) NewListByExts(exts ...string) (*MimeTypeList, error) {
	mimeList := make([]MimeType, 0, len(exts))

	for _, ext := range exts {
		contentType, err := mt.getContentType(ext)
		if err != nil {
			return nil, err
		}

		mimeList = append(
			mimeList,
			MimeType{
				Extension:   ext,
				ContentType: contentType,
			},
		)
	}

	return NewMimeTypeList(mimeList), nil
}

// CheckExt  - comment func.
func (mt *MimeTypeList) CheckExt(ext string) error {
	if _, err := mt.getContentType(ext); err != nil {
		return err
	}

	return nil
}

// CheckExtByFileName  - comment func.
func (mt *MimeTypeList) CheckExtByFileName(name string) error {
	if _, err := mt.getContentType(path.Ext(name)); err != nil {
		return err
	}

	return nil
}

// CheckContentType  - comment func.
func (mt *MimeTypeList) CheckContentType(contentType string) error {
	if _, err := mt.getExt(contentType); err != nil {
		return err
	}

	return nil
}

// ContentType  - comment func.
func (mt *MimeTypeList) ContentType(ext string) string {
	value, err := mt.getContentType(ext)
	if err != nil {
		return ""
	}

	return value
}

// ContentTypeByFileName  - comment func.
func (mt *MimeTypeList) ContentTypeByFileName(name string) string {
	value, err := mt.getContentType(path.Ext(name))
	if err != nil {
		return ""
	}

	return value
}

// Ext  - comment func.
func (mt *MimeTypeList) Ext(contentType string) string {
	value, err := mt.getExt(contentType)
	if err != nil {
		return ""
	}

	return value
}

func (mt *MimeTypeList) getContentType(ext string) (string, error) {
	if ext == "" || len(ext) == 1 && ext[0] == '.' {
		return "", errors.New("arg 'ext' is empty")
	}

	if ext[0] == '.' { // если указано расширение с точкой в начале
		ext = ext[1:]
	}

	if value, ok := mt.extMap[ext]; ok {
		return value, nil
	}

	return "", fmt.Errorf("mime not found for arg '%s'", ext)
}

func (mt *MimeTypeList) getExt(contentType string) (string, error) {
	if contentType == "" {
		return "", errors.New("arg 'contentType' is empty")
	}

	if value, ok := mt.contentTypeMap[contentType]; ok {
		return value, nil
	}

	return "", fmt.Errorf("ext not found for arg '%s'", contentType)
}
