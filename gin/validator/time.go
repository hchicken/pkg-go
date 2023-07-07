package validator

import (
	"strconv"
	"time"

	"github.com/hchicken/pkg-go/util"

	"github.com/go-playground/validator/v10"
)

// timeStringToUnix 检查时间类型
func timeStringToUnix(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v != "" {
		timeUnix, err := util.TimeToUnixV2(v)
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
