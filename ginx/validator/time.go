package validator

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hchicken/pkg-go/date"
)

// timeStringToUnix 检查时间类型
func timeStringToUnix(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v != "" {
		timeUnix, err := date.TimeToUnixV2(v)
		if err != nil {
			return false
		}
		f := fl.Field()
		if f.CanSet() {
			myTime := strconv.FormatInt(timeUnix, 10)
			f.SetString(myTime)
		}
	}
	return true
}

// validTimeString 检查时间类型
func validTimeString(fl validator.FieldLevel) bool {
	v := fl.Field().String()

	if v != "" {
		_, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return false
		}
	}
	return true
}
