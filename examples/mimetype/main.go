package main

import (
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
)

func main() {
	logger := mrlog.Default().With().Str("example", "mimetype").Logger()

	mime := mrlib.NewMimeTypeList(logger, getMimeTypeListFromConfig())
	jpegZip, _ := mime.MimeTypesByExts([]string{".jpeg", ".zip"})

	logger.Info().Msgf(".json content-type: [%s]", mime.ContentType("json"))
	logger.Info().Msgf(".jpeg content-type: [%s]", mime.ContentType(".jpeg"))
	logger.Info().Msgf(".pdf content-type: [%s]", mime.ContentTypeByFileName("file-name.pdf"))
	logger.Info().Msgf("image/jpeg extension: [.%s]", mime.Ext("image/jpeg"))
	logger.Info().Msgf("jpeg/zip mime-types: [.%v]", jpegZip)
}

func getMimeTypeListFromConfig() []mrlib.MimeType {
	return []mrlib.MimeType{
		{ContentType: "application/excel", Extension: "xl"},
		{ContentType: "application/mac-binhex40", Extension: "hqx"},
		{ContentType: "application/mac-compactpro", Extension: "cpt"},
		{ContentType: "application/macbinary", Extension: "bin"},
		{ContentType: "application/msword", Extension: "word"},
		{ContentType: "application/octet-stream", Extension: "class"},
		{ContentType: "application/octet-stream", Extension: "dll"},
		{ContentType: "application/octet-stream", Extension: "dms"},
		{ContentType: "application/octet-stream", Extension: "exe"},
		{ContentType: "application/octet-stream", Extension: "lha"},
		{ContentType: "application/octet-stream", Extension: "lzh"},
		{ContentType: "application/octet-stream", Extension: "psd"},
		{ContentType: "application/octet-stream", Extension: "sea"},
		{ContentType: "application/octet-stream", Extension: "so"},
		{ContentType: "application/oda", Extension: "oda"},
		{ContentType: "application/pdf", Extension: "pdf"},
		{ContentType: "application/postscript", Extension: "ai"},
		{ContentType: "application/postscript", Extension: "eps"},
		{ContentType: "application/postscript", Extension: "ps"},
		{ContentType: "application/smil", Extension: "smi"},
		{ContentType: "application/smil", Extension: "smil"},
		{ContentType: "application/vnd.mif", Extension: "mif"},
		{ContentType: "application/vnd.wap.wbxml", Extension: "wbxml"},
		{ContentType: "application/vnd.wap.wmlc", Extension: "wmlc"},
		{ContentType: "application/x-director", Extension: "dcr"},
		{ContentType: "application/x-director", Extension: "dir"},
		{ContentType: "application/x-director", Extension: "dxr"},
		{ContentType: "application/x-dvi", Extension: "dvi"},
		{ContentType: "application/x-gtar", Extension: "gtar"},
		{ContentType: "application/x-httpd-php", Extension: "php3"},
		{ContentType: "application/x-httpd-php", Extension: "php4"},
		{ContentType: "application/x-httpd-php", Extension: "php"},
		{ContentType: "application/x-httpd-php", Extension: "phtml"},
		{ContentType: "application/x-httpd-php-source", Extension: "phps"},
		{ContentType: "application/x-javascript", Extension: "js"},
		{ContentType: "application/json", Extension: "json"},
		{ContentType: "application/x-shockwave-flash", Extension: "swf"},
		{ContentType: "application/x-stuffit", Extension: "sit"},
		{ContentType: "application/vnd.rar", Extension: "rar"},
		{ContentType: "application/x-tar", Extension: "tar"},
		{ContentType: "application/x-tar", Extension: "tgz"},
		{ContentType: "application/xhtml+xml", Extension: "xht"},
		{ContentType: "application/xhtml+xml", Extension: "xhtml"},
		{ContentType: "application/zip", Extension: "zip"},

		{ContentType: "audio/midi", Extension: "mid"},
		{ContentType: "audio/midi", Extension: "midi"},
		{ContentType: "audio/mpeg", Extension: "mp2"},
		{ContentType: "audio/mpeg", Extension: "mp3"},
		{ContentType: "audio/mpeg", Extension: "mpga"},
		{ContentType: "audio/x-aiff", Extension: "aif"},
		{ContentType: "audio/x-aiff", Extension: "aifc"},
		{ContentType: "audio/x-aiff", Extension: "aiff"},
		{ContentType: "audio/x-pn-realaudio", Extension: "ram"},
		{ContentType: "audio/x-pn-realaudio", Extension: "rm"},
		{ContentType: "audio/x-pn-realaudio-plugin", Extension: "rpm"},
		{ContentType: "audio/x-realaudio", Extension: "ra"},
		{ContentType: "audio/x-wav", Extension: "wav"},

		{ContentType: "image/bmp", Extension: "bmp"},
		{ContentType: "image/gif", Extension: "gif"},
		{ContentType: "image/jpeg", Extension: "jpg"},
		{ContentType: "image/jpeg", Extension: "jpeg"},
		{ContentType: "image/jpeg", Extension: "jpe"},
		{ContentType: "image/png", Extension: "png"},
		{ContentType: "image/tiff", Extension: "tiff"},
		{ContentType: "image/tiff", Extension: "tif"},

		{ContentType: "message/rfc822", Extension: "eml"},
		{ContentType: "text/css", Extension: "css"},
		{ContentType: "text/html", Extension: "html"},
		{ContentType: "text/html", Extension: "htm"},
		{ContentType: "text/html", Extension: "shtml"},
		{ContentType: "text/plain", Extension: "log"},
		{ContentType: "text/plain", Extension: "text"},
		{ContentType: "text/plain", Extension: "txt"},
		{ContentType: "text/richtext", Extension: "rtx"},
		{ContentType: "text/rtf", Extension: "rtf"},
		{ContentType: "text/vcard", Extension: "vcf"},
		{ContentType: "text/vcard", Extension: "vcard"},
		{ContentType: "text/xml", Extension: "xml"},
		{ContentType: "text/xml", Extension: "xsl"},

		{ContentType: "video/mpeg", Extension: "mpg"},
		{ContentType: "video/mpeg", Extension: "mpeg"},
		{ContentType: "video/mpeg", Extension: "mpe"},
		{ContentType: "video/mp4", Extension: "mp4"},
		{ContentType: "video/quicktime", Extension: "mov"},
		{ContentType: "video/quicktime", Extension: "qt"},
		{ContentType: "video/vnd.rn-realvideo", Extension: "rv"},
		{ContentType: "video/x-msvideo", Extension: "avi"},
		{ContentType: "video/x-sgi-movie", Extension: "movie"},

		{ContentType: "application/msword", Extension: "doc"},
		{ContentType: "application/msword", Extension: "dot"},
		{ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document", Extension: "docx"},
		{ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.template", Extension: "dotx"},
		{ContentType: "application/vnd.ms-word.document.macroEnabled.12", Extension: "docm"},
		{ContentType: "application/vnd.ms-word.template.macroEnabled.12", Extension: "dotm"},
		{ContentType: "application/vnd.ms-excel", Extension: "xls"},
		{ContentType: "application/vnd.ms-excel", Extension: "xlt"},
		{ContentType: "application/vnd.ms-excel", Extension: "xla"},
		{ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Extension: "xlsx"},
		{ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.template", Extension: "xltx"},
		{ContentType: "application/vnd.ms-excel.sheet.macroEnabled.12", Extension: "xlsm"},
		{ContentType: "application/vnd.ms-excel.template.macroEnabled.12", Extension: "xltm"},
		{ContentType: "application/vnd.ms-excel.addin.macroEnabled.12", Extension: "xlam"},
		{ContentType: "application/vnd.ms-excel.sheet.binary.macroEnabled.12", Extension: "xlsb"},
		{ContentType: "application/vnd.ms-powerpoint", Extension: "ppt"},
		{ContentType: "application/vnd.ms-powerpoint", Extension: "pot"},
		{ContentType: "application/vnd.ms-powerpoint", Extension: "pps"},
		{ContentType: "application/vnd.ms-powerpoint", Extension: "ppa"},
		{ContentType: "application/vnd.openxmlformats-officedocument.presentationml.presentation", Extension: "pptx"},
		{ContentType: "application/vnd.openxmlformats-officedocument.presentationml.template", Extension: "potx"},
		{ContentType: "application/vnd.openxmlformats-officedocument.presentationml.slideshow", Extension: "ppsx"},
		{ContentType: "application/vnd.ms-powerpoint.addin.macroEnabled.12", Extension: "ppam"},
		{ContentType: "application/vnd.ms-powerpoint.presentation.macroEnabled.12", Extension: "pptm"},
		{ContentType: "application/vnd.ms-powerpoint.template.macroEnabled.12", Extension: "potm"},
		{ContentType: "application/vnd.ms-powerpoint.slideshow.macroEnabled.12", Extension: "ppsm"},
		{ContentType: "application/vnd.ms-access", Extension: "mdb"},
	}
}
