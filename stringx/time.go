package stringx

import "time"

// NowTimeString get the current time string
func NowTimeString() string {
	timeStr := time.Now().Format("20060102150405")
	return timeStr
}
