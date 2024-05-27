package utils

import "encoding/base64"

// BytesToBase64 converts bytes to base64 string.
// This ensures that the resulting string can be safely used in text-based contexts like cookies, URLs, and HTTP headers,
// which may not handle arbitrary binary data properly.
func BytesToBase64(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

// Base64ToBytes converts string to base64
func Base64ToBytes(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}
