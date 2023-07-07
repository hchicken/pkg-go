package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// 去除 string 值空格
func doTrimStringField(fl validator.FieldLevel) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	v := fl.Field().String()
	f := fl.Field()

	f.SetString(strings.TrimSpace(v))
	return true
}
