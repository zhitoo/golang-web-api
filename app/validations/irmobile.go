package validations

import (
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func IranianMobileNumberValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	if value == "" {
		return true
	}

	matched, err := regexp.MatchString(`^09[0-9]{9}$`, value)
	if err != nil {
		log.Println(err)
		return false
	}
	return matched
}
