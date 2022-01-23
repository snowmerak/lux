package util

func GetContentTypeFromExt(ext string) string {
	contentType := "text/plain"
	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".json":
		contentType = "application/json"
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpeg"
	case ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	case ".svg":
		contentType = "image/svg+xml"
	case ".ico":
		contentType = "image/x-icon"
	case ".woff":
		contentType = "application/font-woff"
	case ".woff2":
		contentType = "application/font-woff2"
	case ".ttf":
		contentType = "application/font-ttf"
	case ".otf":
		contentType = "application/font-otf"
	case ".eot":
		contentType = "application/vnd.ms-fontobject"
	case ".mp4":
		contentType = "video/mp4"
	case ".webm":
		contentType = "video/webm"
	case ".ogv":
		contentType = "video/ogg"
	case ".mp3":
		contentType = "audio/mpeg"
	case ".wav":
		contentType = "audio/wav"
	case ".ogg":
		contentType = "audio/ogg"
	case ".flac":
		contentType = "audio/flac"
	case ".wma":
		contentType = "audio/x-ms-wma"
	case ".aac":
		contentType = "audio/aac"
	case ".m4a":
		contentType = "audio/m4a"
	case ".mpg":
		contentType = "video/mpeg"
	case ".mpeg":
		contentType = "video/mpeg"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".mov":
		contentType = "video/quicktime"
	case ".zip":
		contentType = "application/zip"
	case ".rar":
		contentType = "application/x-rar-compressed"
	case ".7z":
		contentType = "application/x-7z-compressed"
	case ".tar":
		contentType = "application/x-tar"
	case ".gz":
		contentType = "application/gzip"
	case ".bz2":
		contentType = "application/x-bzip2"
	case ".doc":
		contentType = "application/msword"
	case ".docx":
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		contentType = "application/vnd.ms-excel"
	case ".xlsx":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		contentType = "application/vnd.ms-powerpoint"
	case ".pptx":
		contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".pdf":
		contentType = "application/pdf"
	case ".txt":
		contentType = "text/plain"
	case ".rtf":
		contentType = "application/rtf"
	case ".xml":
		contentType = "text/xml"
	case ".xsl":
		contentType = "text/xsl"
	case ".csv":
		contentType = "text/csv"
	case ".tsv":
		contentType = "text/tab-separated-values"
	}
	return contentType
}
