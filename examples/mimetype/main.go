package main

import (
	"log"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/mrlogbase"
)

func main() {
	logger := mrlogbase.New(mrlog.DebugLevel).With().Str("example", "mimetype").Logger()
	if err := mrlog.SetDefault(logger); err != nil {
		log.Fatal(err)
	}

	mime := mrlib.NewMimeTypeList(getMimeTypeListFromConfig())

	logger.Info().Msgf(".json content-type: [%s]", mime.ContentType("json"))
	logger.Info().Msgf(".jpeg content-type: [%s]", mime.ContentType(".jpeg"))
	logger.Info().Msgf(".pdf content-type: [%s]", mime.ContentTypeByFileName("file-name.pdf"))
	logger.Info().Msgf("image/jpeg extension: [.%s]", mime.Ext("image/jpeg"))
}

func getMimeTypeListFromConfig() []mrlib.MimeType {
	return []mrlib.MimeType{
		{Extension: "xl", ContentType: "application/excel"},
		{Extension: "hqx", ContentType: "application/mac-binhex40"},
		{Extension: "cpt", ContentType: "application/mac-compactpro"},
		{Extension: "bin", ContentType: "application/macbinary"},
		{Extension: "word", ContentType: "application/msword"},
		{Extension: "class", ContentType: "application/octet-stream"},
		{Extension: "dll", ContentType: "application/octet-stream"},
		{Extension: "dms", ContentType: "application/octet-stream"},
		{Extension: "exe", ContentType: "application/octet-stream"},
		{Extension: "lha", ContentType: "application/octet-stream"},
		{Extension: "lzh", ContentType: "application/octet-stream"},
		{Extension: "psd", ContentType: "application/octet-stream"},
		{Extension: "sea", ContentType: "application/octet-stream"},
		{Extension: "so", ContentType: "application/octet-stream"},
		{Extension: "oda", ContentType: "application/oda"},
		{Extension: "pdf", ContentType: "application/pdf"},
		{Extension: "ai", ContentType: "application/postscript"},
		{Extension: "eps", ContentType: "application/postscript"},
		{Extension: "ps", ContentType: "application/postscript"},
		{Extension: "smi", ContentType: "application/smil"},
		{Extension: "smil", ContentType: "application/smil"},
		{Extension: "mif", ContentType: "application/vnd.mif"},
		{Extension: "wbxml", ContentType: "application/vnd.wap.wbxml"},
		{Extension: "wmlc", ContentType: "application/vnd.wap.wmlc"},
		{Extension: "dcr", ContentType: "application/x-director"},
		{Extension: "dir", ContentType: "application/x-director"},
		{Extension: "dxr", ContentType: "application/x-director"},
		{Extension: "dvi", ContentType: "application/x-dvi"},
		{Extension: "gtar", ContentType: "application/x-gtar"},
		{Extension: "php3", ContentType: "application/x-httpd-php"},
		{Extension: "php4", ContentType: "application/x-httpd-php"},
		{Extension: "php", ContentType: "application/x-httpd-php"},
		{Extension: "phtml", ContentType: "application/x-httpd-php"},
		{Extension: "phps", ContentType: "application/x-httpd-php-source"},
		{Extension: "js", ContentType: "application/x-javascript"},
		{Extension: "json", ContentType: "application/json"},
		{Extension: "swf", ContentType: "application/x-shockwave-flash"},
		{Extension: "sit", ContentType: "application/x-stuffit"},
		{Extension: "rar", ContentType: "application/vnd.rar"},
		{Extension: "tar", ContentType: "application/x-tar"},
		{Extension: "tgz", ContentType: "application/x-tar"},
		{Extension: "xht", ContentType: "application/xhtml+xml"},
		{Extension: "xhtml", ContentType: "application/xhtml+xml"},
		{Extension: "zip", ContentType: "application/zip"},

		{Extension: "mid", ContentType: "audio/midi"},
		{Extension: "midi", ContentType: "audio/midi"},
		{Extension: "mp2", ContentType: "audio/mpeg"},
		{Extension: "mp3", ContentType: "audio/mpeg"},
		{Extension: "mpga", ContentType: "audio/mpeg"},
		{Extension: "aif", ContentType: "audio/x-aiff"},
		{Extension: "aifc", ContentType: "audio/x-aiff"},
		{Extension: "aiff", ContentType: "audio/x-aiff"},
		{Extension: "ram", ContentType: "audio/x-pn-realaudio"},
		{Extension: "rm", ContentType: "audio/x-pn-realaudio"},
		{Extension: "rpm", ContentType: "audio/x-pn-realaudio-plugin"},
		{Extension: "ra", ContentType: "audio/x-realaudio"},
		{Extension: "wav", ContentType: "audio/x-wav"},

		{Extension: "bmp", ContentType: "image/bmp"},
		{Extension: "gif", ContentType: "image/gif"},
		{Extension: "jpg", ContentType: "image/jpeg"},
		{Extension: "jpeg", ContentType: "image/jpeg"},
		{Extension: "jpe", ContentType: "image/jpeg"},
		{Extension: "png", ContentType: "image/png"},
		{Extension: "tiff", ContentType: "image/tiff"},
		{Extension: "tif", ContentType: "image/tiff"},

		{Extension: "eml", ContentType: "message/rfc822"},
		{Extension: "css", ContentType: "text/css"},
		{Extension: "html", ContentType: "text/html"},
		{Extension: "htm", ContentType: "text/html"},
		{Extension: "shtml", ContentType: "text/html"},
		{Extension: "log", ContentType: "text/plain"},
		{Extension: "text", ContentType: "text/plain"},
		{Extension: "txt", ContentType: "text/plain"},
		{Extension: "rtx", ContentType: "text/richtext"},
		{Extension: "rtf", ContentType: "text/rtf"},
		{Extension: "vcf", ContentType: "text/vcard"},
		{Extension: "vcard", ContentType: "text/vcard"},
		{Extension: "xml", ContentType: "text/xml"},
		{Extension: "xsl", ContentType: "text/xml"},

		{Extension: "mpg", ContentType: "video/mpeg"},
		{Extension: "mpeg", ContentType: "video/mpeg"},
		{Extension: "mpe", ContentType: "video/mpeg"},
		{Extension: "mp4", ContentType: "video/mp4"},
		{Extension: "mov", ContentType: "video/quicktime"},
		{Extension: "qt", ContentType: "video/quicktime"},
		{Extension: "rv", ContentType: "video/vnd.rn-realvideo"},
		{Extension: "avi", ContentType: "video/x-msvideo"},
		{Extension: "movie", ContentType: "video/x-sgi-movie"},

		{Extension: "doc", ContentType: "application/msword"},
		{Extension: "dot", ContentType: "application/msword"},
		{Extension: "docx", ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		{Extension: "dotx", ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.template"},
		{Extension: "docm", ContentType: "application/vnd.ms-word.document.macroEnabled.12"},
		{Extension: "dotm", ContentType: "application/vnd.ms-word.template.macroEnabled.12"},
		{Extension: "xls", ContentType: "application/vnd.ms-excel"},
		{Extension: "xlt", ContentType: "application/vnd.ms-excel"},
		{Extension: "xla", ContentType: "application/vnd.ms-excel"},
		{Extension: "xlsx", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		{Extension: "xltx", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.template"},
		{Extension: "xlsm", ContentType: "application/vnd.ms-excel.sheet.macroEnabled.12"},
		{Extension: "xltm", ContentType: "application/vnd.ms-excel.template.macroEnabled.12"},
		{Extension: "xlam", ContentType: "application/vnd.ms-excel.addin.macroEnabled.12"},
		{Extension: "xlsb", ContentType: "application/vnd.ms-excel.sheet.binary.macroEnabled.12"},
		{Extension: "ppt", ContentType: "application/vnd.ms-powerpoint"},
		{Extension: "pot", ContentType: "application/vnd.ms-powerpoint"},
		{Extension: "pps", ContentType: "application/vnd.ms-powerpoint"},
		{Extension: "ppa", ContentType: "application/vnd.ms-powerpoint"},
		{Extension: "pptx", ContentType: "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		{Extension: "potx", ContentType: "application/vnd.openxmlformats-officedocument.presentationml.template"},
		{Extension: "ppsx", ContentType: "application/vnd.openxmlformats-officedocument.presentationml.slideshow"},
		{Extension: "ppam", ContentType: "application/vnd.ms-powerpoint.addin.macroEnabled.12"},
		{Extension: "pptm", ContentType: "application/vnd.ms-powerpoint.presentation.macroEnabled.12"},
		{Extension: "potm", ContentType: "application/vnd.ms-powerpoint.template.macroEnabled.12"},
		{Extension: "ppsm", ContentType: "application/vnd.ms-powerpoint.slideshow.macroEnabled.12"},
		{Extension: "mdb", ContentType: "application/vnd.ms-access"},
	}
}
