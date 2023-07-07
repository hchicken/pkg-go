package util

import (
	"time"

	"github.com/araddon/dateparse"
)

// UnixToTime 时间戳转时间
func UnixToTime(unix int64) string {
	return time.Unix(unix, 0).Format("2006-01-02 15:04:05") // 时间戳转时间字符串
}

// TimeToUnix 时间字符串转时间戳
func TimeToUnix(str string) (int64, error) {
	the_time, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		return 0, err
	}
	unix := the_time.Unix()
	return unix, nil
}

// TimeToUnix1 自定义事件格式
func TimeToUnix1(layout, value string) (int64, error) {
	theTime, err := time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		return 0, err
	}
	unix := theTime.Unix()
	return unix, nil
}

// TimeToUnixV2 兼容各种时间格式
func TimeToUnixV2(value string) (int64, error) {
	t1, err := dateparse.ParseAny(value)
	if err != nil {
		return 0, err
	}
	unix := t1.Unix()
	return unix, nil
}
