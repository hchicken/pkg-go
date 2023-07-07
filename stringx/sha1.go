package stringx

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

// Sha1Base64 sha1加密之后用base64处理
func Sha1Base64(key, sig string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(sig))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Sha1 sha1加密
func Sha1(key, sig string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(sig))
	return hex.EncodeToString(mac.Sum(nil))
}
