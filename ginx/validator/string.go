package validator

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Trim spaces from string values
func doTrimStringField(fl validator.FieldLevel) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error in doTrimStringField: %v", err)
		}
	}()

	v := fl.Field().String()
	f := fl.Field()

	f.SetString(strings.TrimSpace(v))
	return true
}
