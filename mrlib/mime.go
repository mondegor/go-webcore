package mrlib

import (
	"errors"
	"fmt"
	"strings"
)

type (
	// MimeTypeList - хранит соответствие расширений их типам файлов (в обе стороны).
	MimeTypeList struct {
		contentTypeMap map[string]string
		extensionMap   map[string]string
	}

	// MimeType - хранит расширение и соответствующий ему тип файла.
	MimeType struct {
		ContentType string `yaml:"type"`
		Extension   string `yaml:"ext"`
	}
)

// NewMimeTypeList - создаёт объект MimeTypeList на основе списка соответствий расширений и файлов.
func NewMimeTypeList(items []MimeType) *MimeTypeList {
	mimeMap := make(map[string]string, len(items))
	extMap := make(map[string]string, len(items))

	for _, item := range items {
		item.Extension = strings.TrimPrefix(item.Extension, ".")
		extMap[item.Extension] = item.ContentType

		// т.к. у одного типа может быть несколько расширений,
		// то в индекс попадает только первый зарегистрированный
		if _, ok := mimeMap[item.ContentType]; !ok {
			mimeMap[item.ContentType] = item.Extension
		}
	}

	return &MimeTypeList{
		contentTypeMap: mimeMap,
		extensionMap:   extMap,
	}
}

// MimeTypesByExts - возвращает MimeType массив, в который войдут указанные расширения,
// если хотя бы одно расширение не зарегистрировано в списке, то будет выдана ошибка.
func (mt *MimeTypeList) MimeTypesByExts(values []string) ([]MimeType, error) {
	mime := make([]MimeType, len(values))

	for i, ext := range values {
		contentType, err := mt.ContentTypeByExt(ext)
		if err != nil {
			return nil, err
		}

		mime[i] = MimeType{
			ContentType: contentType,
			Extension:   ext,
		}
	}

	return mime, nil
}

// ContentTypeByExt - возвращает тип файла по указанному расширению,
// если тип не найден, то возвращается пустая строка.
func (mt *MimeTypeList) ContentTypeByExt(value string) (string, error) {
	if value == "" || len(value) == 1 && value[0] == '.' {
		return "", errors.New("arg 'value' is empty")
	}

	if value[0] == '.' { // если указано расширение с точкой в начале
		value = value[1:]
	}

	if ext, ok := mt.extensionMap[value]; ok {
		return ext, nil
	}

	return "", fmt.Errorf("mime not found for arg '%s'", value)
}

// ExtByContentType - возвращает расширение по указанному типу файла,
// если расширение не найдено, то возвращается пустая строка.
func (mt *MimeTypeList) ExtByContentType(value string) (string, error) {
	if value == "" {
		return "", errors.New("arg 'value' is empty")
	}

	if contentType, ok := mt.contentTypeMap[value]; ok {
		return contentType, nil
	}

	return "", fmt.Errorf("ext not found for arg '%s'", value)
}
