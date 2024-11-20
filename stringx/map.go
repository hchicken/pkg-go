package stringx

import "strings"

// ParseKeyValuePairs 解析key=value;key1=value2;
func ParseKeyValuePairs(input string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(strings.TrimSuffix(input, ";"), ";")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return result
}
