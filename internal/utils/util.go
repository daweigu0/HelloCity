package utils

// GetFileType 根据文件的扩展名获取文件对应的类型
func GetFileType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp", ".svg", ".heif", ".heic":
		return "image"
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm", ".m4v", ".3gp", ".3gpp", ".mpeg", ".mpg":
		return "video"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma", ".m4a", ".opus", ".aiff", ".aif":
		return "audio"
	default:
		return "unknown"
	}
}
