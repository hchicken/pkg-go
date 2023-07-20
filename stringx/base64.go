package stringx

import (
	"encoding/base64"
	"os"
)

// Base64Encode encodes input bytes to base64
func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

// StrToBase64 字符串转base64
func StrToBase64(data string) string {
	return Base64Encode([]byte(data))
}

// FileToBase64 文件转base64
func FileToBase64(file string) (string, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return Base64Encode(f), nil
}
