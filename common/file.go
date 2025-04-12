package common

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"
)

// IsValidImage -
func IsValidImage(fileHeader *multipart.FileHeader) (bool, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return false, "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		return false, "", err
	}

	if _, err = file.Seek(0, 0); err != nil {
		return false, "", err
	}

	// 检查MIME类型
	mimeType := http.DetectContentType(buffer)
	if !strings.HasPrefix(mimeType, "image/") {
		return false, "", nil
	}

	// 检查常见图片格式签名
	switch {
	case bytes.HasPrefix(buffer, []byte("\xFF\xD8\xFF")):
		return true, "jpeg", nil
	case bytes.HasPrefix(buffer, []byte("\x89PNG\r\n\x1a\n")):
		return true, "png", nil
	case bytes.HasPrefix(buffer, []byte("GIF87a")), bytes.HasPrefix(buffer, []byte("GIF89a")):
		return true, "gif", nil
	case bytes.HasPrefix(buffer, []byte("BM")):
		return true, "bmp", nil
	case len(buffer) > 12 && bytes.Equal(buffer[8:12], []byte("WEBP")):
		return true, "webp", nil
	default:
		return false, "", nil
	}
}
