package token

import "encoding/base64"

// Encode byte
func encodeToBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Decode string
func decodeBase64String(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
