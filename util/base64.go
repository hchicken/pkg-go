package util

import (
	"encoding/base64"
	"os"
)

// StrToBase64 字符串转base64
func StrToBase64(data string) (string, error) {
	sDec := base64.StdEncoding.EncodeToString([]byte(data))
	return sDec, nil
}

// FileToBase64 文件转base64
func FileToBase64(file string) (string, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	sDec := base64.StdEncoding.EncodeToString(f)
	return sDec, nil
}
