package util

import "os"

// GetEnv 默认会设置env
func GetEnv(key, value string) string {
	v := os.Getenv(key)
	if v == "" {
		v = value
	}
	return v
}
