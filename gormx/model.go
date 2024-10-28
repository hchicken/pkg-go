package gormx

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// TabBaseModel 模型定义基类
type TabBaseModel struct {
	CreatedBy string   `json:"created_by" gorm:"type:varchar(128);column:created_by;comment:'添加人'"`
	UpdatedBy string   `json:"updated_by" gorm:"type:varchar(128);column:updated_by;comment:'更新人'"`
	CreatedAt JsonTime `json:"created_at" gorm:"comment:'添加时间'"`
	UpdatedAt JsonTime `json:"updated_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP on update current_timestamp;comment:'更新时间'"`
}

// JsonTime 自定义时间格式(用来处理字符串时间格式)
type JsonTime struct {
	time.Time
}

// MarshalJSON 序列化json
func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(time.DateTime))
	return []byte(formatted), nil
}

// UnmarshalJSON 反序列号json
func (t *JsonTime) UnmarshalJSON(b []byte) error {
	now, err := time.ParseInLocation(`"`+time.DateTime+`"`, string(b), time.Local)
	t.Time = now
	return err
}

// Value ...
func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan ...
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
	return t.Format(time.DateTime)
}
