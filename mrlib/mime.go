package mrlib

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

type (
	// MimeTypeList - хранит соответствие расширений их типам файлов (в обе стороны).
	MimeTypeList struct {
		contentTypeMap map[string]string
		extensionMap   map[string]string
		logger         logger
	}

	// MimeType - хранит расширение и соответствующий ему тип файла.
	MimeType struct {
		ContentType string `yaml:"type"`
		Extension   string `yaml:"ext"`
	}
)

// NewMimeTypeList - создаёт объект MimeTypeList на основе списка соответствий расширений и файлов.
func NewMimeTypeList(logger logger, items []MimeType) *MimeTypeList {
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
		logger:         logger,
	}
}

// MimeTypesByExts - возвращает MimeType массив, в который войдут указанные расширения,
// если хотя бы одно расширение не зарегистрировано в списке, то будет выдана ошибка.
func (mt *MimeTypeList) MimeTypesByExts(values []string) ([]MimeType, error) {
	mime := make([]MimeType, len(values))

	for i, ext := range values {
		contentType, err := mt.getContentType(ext)
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

// CheckExt - возвращает ошибку, если указанное расширение не зарегистрировано в списке.
func (mt *MimeTypeList) CheckExt(ext string) error {
	if _, err := mt.getContentType(ext); err != nil {
		return err
	}

	return nil
}

// CheckExtByFileName - возвращает ошибку, если расширение указанного файла не зарегистрировано в списке.
func (mt *MimeTypeList) CheckExtByFileName(name string) error {
	if _, err := mt.getContentType(path.Ext(name)); err != nil {
		return err
	}

	return nil
}

// CheckContentType - возвращает ошибку, если указанный тип файла не зарегистрирован в списке.
func (mt *MimeTypeList) CheckContentType(contentType string) error {
	if _, err := mt.getExt(contentType); err != nil {
		return err
	}

	return nil
}

// ContentType - возвращает тип файла по указанному расширению,
// если тип не найден, то возвращается пустая строка.
func (mt *MimeTypeList) ContentType(ext string) string {
	value, err := mt.getContentType(ext)
	if err != nil {
		mt.logger.Printf(fmt.Errorf("ContentType: %w", err).Error())

		return ""
	}

	return value
}

// ContentTypeByFileName - возвращает тип файла по расширению указанного файла,
// если тип не найден, то возвращается пустая строка.
func (mt *MimeTypeList) ContentTypeByFileName(name string) string {
	value, err := mt.getContentType(path.Ext(name))
	if err != nil {
		mt.logger.Printf(fmt.Errorf("ContentTypeByFileName: %w", err).Error())

		return ""
	}

	return value
}

// Ext - возвращает расширение по указанному типу файла,
// если расширение не найдено, то возвращается пустая строка.
func (mt *MimeTypeList) Ext(contentType string) string {
	value, err := mt.getExt(contentType)
	if err != nil {
		mt.logger.Printf(fmt.Errorf("Ext: %w", err).Error())

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

	if value, ok := mt.extensionMap[ext]; ok {
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
