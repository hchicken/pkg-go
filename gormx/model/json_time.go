package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JsonTime 自定义时间格式(用来处理字符串时间格式)
type JsonTime struct {
	time.Time
}

const baseTimeFormat string = "2006-01-02 15:04:05"

// MarshalJSON 序列化json
func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(baseTimeFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON 反序列号json
func (t *JsonTime) UnmarshalJSON(b []byte) error {
	now, err := time.ParseInLocation(`"`+baseTimeFormat+`"`, string(b), time.Local)
	t.Time = now
	return err
}

// Value TODO
// value
func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan TODO
// scan
func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to JSONTime", v)
}

// FormatTime 时间转换
func (t *JsonTime) FormatTime() string {
	return t.Format("2006-01-02 15:04:05")
}
