package stringx

import (
	"strings"

	jsonIterator "github.com/json-iterator/go"
	"github.com/valyala/fasttemplate"
)

// Format 格式化字符串
func Format(s, start, end string, m map[string]interface{}) (string, error) {
	newM := make(map[string]interface{}, len(m))
	json := jsonIterator.ConfigCompatibleWithStandardLibrary // 创建jsonIterator实例
	for k, v := range m {
		b, err := json.MarshalToString(v)
		if err != nil {
			return "", err
		}
		newM[k] = strings.Trim(b, "\"")
	}
	value := fasttemplate.New(s, start, end).ExecuteString(newM)
	return value, nil
}
