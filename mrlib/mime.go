package mrlib

import (
	"path"
)

var (
	mimeExtType = map[string]string{
		// "xl": "application/excel",
		// "hqx": "application/mac-binhex40",
		// "cpt": "application/mac-compactpro",
		// "bin": "application/macbinary",
		// "word": "application/msword",
		// "class": "application/octet-stream",
		// "dll": "application/octet-stream",
		// "dms": "application/octet-stream",
		// "exe": "application/octet-stream",
		// "lha": "application/octet-stream",
		// "lzh": "application/octet-stream",
		// "psd": "application/octet-stream",
		// "sea": "application/octet-stream",
		// "so": "application/octet-stream",
		// "oda": "application/oda",
		"pdf": "application/pdf",
		// "ai": "application/postscript",
		// "eps": "application/postscript",
		// "ps": "application/postscript",
		// "smi": "application/smil",
		// "smil": "application/smil",
		// "mif": "application/vnd.mif",
		// "wbxml": "application/vnd.wap.wbxml",
		// "wmlc": "application/vnd.wap.wmlc",
		// "dcr": "application/x-director",
		// "dir": "application/x-director",
		// "dxr": "application/x-director",
		// "dvi": "application/x-dvi",
		// "gtar": "application/x-gtar",
		// "php3": "application/x-httpd-php",
		// "php4": "application/x-httpd-php",
		// "php": "application/x-httpd-php",
		// "phtml": "application/x-httpd-php",
		// "phps": "application/x-httpd-php-source",
		// "js": "application/x-javascript",
		"json": "application/json",
		// "swf": "application/x-shockwave-flash",
		// "sit": "application/x-stuffit",
		"rar": "application/vnd.rar",
		"tar": "application/x-tar",
		"tgz": "application/x-tar",
		// "xht": "application/xhtml+xml",
		// "xhtml": "application/xhtml+xml",
		"zip": "application/zip",
		// "mid": "audio/midi",
		// "midi": "audio/midi",
		// "mp2": "audio/mpeg",
		// "mp3": "audio/mpeg",
		// "mpga": "audio/mpeg",
		// "aif": "audio/x-aiff",
		// "aifc": "audio/x-aiff",
		// "aiff": "audio/x-aiff",
		// "ram": "audio/x-pn-realaudio",
		// "rm": "audio/x-pn-realaudio",
		// "rpm": "audio/x-pn-realaudio-plugin",
		// "ra": "audio/x-realaudio",
		// "wav": "audio/x-wav",
		// "bmp": "image/bmp",
		"gif":  "image/gif",
		"jpeg": "image/jpeg",
		"jpe":  "image/jpeg",
		"jpg":  "image/jpeg",
		"png":  "image/png",
		// "tiff": "image/tiff",
		// "tif":  "image/tiff",
		// "eml": "message/rfc822",
		"css":  "text/css",
		"html": "text/html",
		// "htm": "text/html",
		// "shtml": "text/html",
		// "log": "text/plain",
		// "text": "text/plain",
		// "txt": "text/plain",
		// "rtx": "text/richtext",
		// "rtf": "text/rtf",
		// "vcf": "text/vcard",
		// "vcard": "text/vcard",
		"xml": "text/xml",
		"xsl": "text/xml",
		// "mpeg": "video/mpeg",
		// "mpe": "video/mpeg",
		// "mpg": "video/mpeg",
		// "mov": "video/quicktime",
		// "qt": "video/quicktime",
		// "rv": "video/vnd.rn-realvideo",
		// "avi": "video/x-msvideo",
		// "movie": "video/x-sgi-movie",

		"doc": "application/msword",
		// "dot":   "application/msword",
		"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		// "dotx":   "application/vnd.openxmlformats-officedocument.wordprocessingml.template",
		// "docm":   "application/vnd.ms-word.document.macroEnabled.12",
		// "dotm":   "application/vnd.ms-word.template.macroEnabled.12",
		"xls": "application/vnd.ms-excel",
		// "xlt":   "application/vnd.ms-excel",
		// "xla":   "application/vnd.ms-excel",
		"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		// "xltx":   "application/vnd.openxmlformats-officedocument.spreadsheetml.template",
		// "xlsm":   "application/vnd.ms-excel.sheet.macroEnabled.12",
		// "xltm":   "application/vnd.ms-excel.template.macroEnabled.12",
		// "xlam":   "application/vnd.ms-excel.addin.macroEnabled.12",
		// "xlsb":   "application/vnd.ms-excel.sheet.binary.macroEnabled.12",
		"ppt": "application/vnd.ms-powerpoint",
		// "pot":   "application/vnd.ms-powerpoint",
		// "pps":   "application/vnd.ms-powerpoint",
		// "ppa":   "application/vnd.ms-powerpoint",
		"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		// "potx":   "application/vnd.openxmlformats-officedocument.presentationml.template",
		// "ppsx":   "application/vnd.openxmlformats-officedocument.presentationml.slideshow",
		// "ppam":   "application/vnd.ms-powerpoint.addin.macroEnabled.12",
		// "pptm":   "application/vnd.ms-powerpoint.presentation.macroEnabled.12",
		// "potm":   "application/vnd.ms-powerpoint.template.macroEnabled.12",
		// "ppsm":   "application/vnd.ms-powerpoint.slideshow.macroEnabled.12",
		// "mdb":   "application/vnd.ms-access",
	}
)

func MimeTypeByExt(ext string) string {
	if ext != "" {
		if ext[0] == '.' {
			ext = ext[1:]
		}

		if value, ok := mimeExtType[ext]; ok {
			return value
		}
	}

	return ""
}

func MimeTypeByFile(name string) string {
	return MimeTypeByExt(path.Ext(name))
}

// MimeType - возвращает value если оно не пустое иначе вычисляется тип по расширению файла
func MimeType(value, fileName string) string {
	if value != "" || fileName == "" {
		return value
	}

	return MimeTypeByFile(fileName)
}
