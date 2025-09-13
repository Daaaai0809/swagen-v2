package constants

const (
	TEXT_PLAIN_ID = iota
	TEXT_HTML_ID
	TEXT_XML_ID
	TEXT_CSS_ID
	TEXT_JAVASCRIPT_ID
	TEXT_CSV_ID
	TEXT_EVENT_STREAM_ID
	APPLICATION_JSON_ID
	APPLICATION_JSON_PATCH_JSON_ID
	APPLICATION_LD_JSON_ID
	APPLICATION_VND_API_JSON_ID
	APPLICATION_XML_ID
	APPLICATION_ATOM_XML_ID
	APPLICATION_RSS_XML_ID
	APPLICATION_OCTET_STREAM_ID
	APPLICATION_PDF_ID
	APPLICATION_ZIP_ID
	APPLICATION_X_WWW_FORM_URLENCODED_ID
	APPLICATION_X_HTML_XML_ID
	APPLICATION_VND_MS_EXCEL_ID
	APPLICATION_VND_MS_POWERPOINT_ID
	APPLICATION_MS_WORD_ID
	MULTIPART_FORM_DATA_ID
	MULTIPART_BYTERANGES_ID
	IMAGE_PNG_ID
	IMAGE_JPEG_ID
	IMAGE_GIF_ID
	IMAGE_SVG_XML_ID
	AUDIO_MPEG_ID
	AUDIO_OGG_ID
	AUDIO_WAV_ID
	VIDEO_MP4_ID
	VIDEO_OGG_ID
	VIDEO_WEBM_ID

	TEXT_PLAIN                        = "text/plain"
	TEXT_HTML                         = "text/html"
	TEXT_XML                          = "text/xml"
	TEXT_CSS                          = "text/css"
	TEXT_JAVASCRIPT                   = "text/javascript"
	TEXT_CSV                          = "text/csv"
	TEXT_EVENT_STREAM                 = "text/event-stream"
	APPLICATION_JSON                  = "application/json"
	APPLICATION_JSON_PATCH_JSON       = "application/json-patch+json"
	APPLICATION_LD_JSON               = "application/ld+json"
	APPLICATION_VND_API_JSON          = "application/vnd.api+json"
	APPLICATION_XML                   = "application/xml"
	APPLICATION_ATOM_XML              = "application/atom+xml"
	APPLICATION_RSS_XML               = "application/rss+xml"
	APPLICATION_OCTET_STREAM          = "application/octet-stream"
	APPLICATION_PDF                   = "application/pdf"
	APPLICATION_ZIP                   = "application/zip"
	APPLICATION_X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded"
	APPLICATION_X_HTML_XML            = "application/xhtml+xml"
	APPLICATION_VND_MS_EXCEL          = "application/vnd.ms-excel"
	APPLICATION_VND_MS_POWERPOINT     = "application/vnd.ms-powerpoint"
	APPLICATION_MS_WORD               = "application/msword"
	MULTIPART_FORM_DATA               = "multipart/form-data"
	MULTIPART_BYTERANGES              = "multipart/byteranges"
	IMAGE_PNG                         = "image/png"
	IMAGE_JPEG                        = "image/jpeg"
	IMAGE_GIF                         = "image/gif"
	IMAGE_SVG_XML                     = "image/svg+xml"
	AUDIO_MPEG                        = "audio/mpeg"
	AUDIO_OGG                         = "audio/ogg"
	AUDIO_WAV                         = "audio/wav"
	VIDEO_MP4                         = "video/mp4"
	VIDEO_OGG                         = "video/ogg"
	VIDEO_WEBM                        = "video/webm"

	_MIME_KEY_JSON = "json"
	_MIME_KEY_HTML = "html"
	_MIME_KEY_XML  = "xml"
	_MIME_KEY_CSV  = "csv"
	_MIME_KEY_FORM = "form"
	_MIME_KEY_TEXT = "text"
	_MIME_KEY_BIN  = "bin"
	_MIME_KEY_SSE  = "sse"
)

var MediaTypeList = []string{
	TEXT_PLAIN,
	TEXT_HTML,
	TEXT_XML,
	TEXT_CSS,
	TEXT_JAVASCRIPT,
	TEXT_CSV,
	APPLICATION_JSON,
	APPLICATION_JSON_PATCH_JSON,
	APPLICATION_LD_JSON,
	APPLICATION_VND_API_JSON,
	APPLICATION_XML,
	APPLICATION_ATOM_XML,
	APPLICATION_RSS_XML,
	APPLICATION_OCTET_STREAM,
	APPLICATION_PDF,
	APPLICATION_ZIP,
	APPLICATION_X_WWW_FORM_URLENCODED,
	APPLICATION_X_HTML_XML,
	APPLICATION_VND_MS_EXCEL,
	APPLICATION_VND_MS_POWERPOINT,
	APPLICATION_MS_WORD,
	MULTIPART_FORM_DATA,
	MULTIPART_BYTERANGES,
	IMAGE_PNG,
	IMAGE_JPEG,
	IMAGE_GIF,
	IMAGE_SVG_XML,
	AUDIO_MPEG,
	AUDIO_OGG,
	AUDIO_WAV,
	VIDEO_MP4,
	VIDEO_OGG,
	VIDEO_WEBM,
}

var MimeKeys = []string{
	_MIME_KEY_JSON,
	_MIME_KEY_HTML,
	_MIME_KEY_XML,
	_MIME_KEY_CSV,
	_MIME_KEY_FORM,
	_MIME_KEY_TEXT,
	_MIME_KEY_BIN,
	_MIME_KEY_SSE,
}

type MinimumMediaTypeMap map[string]string

var MediaTypeMap = MinimumMediaTypeMap{
	_MIME_KEY_JSON: APPLICATION_JSON,
	_MIME_KEY_HTML: TEXT_HTML,
	_MIME_KEY_XML:  APPLICATION_XML,
	_MIME_KEY_CSV:  TEXT_CSV,
	_MIME_KEY_FORM: MULTIPART_FORM_DATA,
	_MIME_KEY_TEXT: TEXT_PLAIN,
	_MIME_KEY_BIN:  APPLICATION_OCTET_STREAM,
	_MIME_KEY_SSE:  TEXT_EVENT_STREAM,
}
