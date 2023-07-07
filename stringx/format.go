package stringx

import (
	"encoding/json"
	"strings"

	"github.com/valyala/fasttemplate"
)

// Format ...
func Format(s, start, end string, m map[string]any) (string, error) {
	newM := make(map[string]any, len(m))
	for k, v := range m {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		newM[k] = strings.ReplaceAll(string(b), "\"", "")
	}
	value := fasttemplate.New(s, start, end).ExecuteString(newM)
	return value, nil
}
