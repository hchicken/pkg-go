package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitValidator 初始化验证
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("timeUnix", timeStringToUnix)
		_ = v.RegisterValidation("trim", doTrimStringField)
		_ = v.RegisterValidation("timeString", validTimeString)
		_ = v.RegisterValidation("notBlank", notBlank)
	}
}
