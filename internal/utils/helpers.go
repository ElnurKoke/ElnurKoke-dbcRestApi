package utils

import (
	"mime"
	"path/filepath"
)

func GetMimeType(filename string) string {
	ext := filepath.Ext(filename)
	return mime.TypeByExtension(ext)
}

func IsValidMimeType(mimeType string, validTypes []string) bool {
	for _, t := range validTypes {
		if t == mimeType {
			return true
		}
	}
	return false
}
