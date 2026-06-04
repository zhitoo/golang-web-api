package validations

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/zhitoo/golang-web-api/lang"
)

type ValidationError struct {
	Property string `json:"property"`
	Tag      string `json:"tag"`
	Value    string `json:"value"`
	Message  string `json:"message"`
}

func GetValidationErrors(err error) *[]ValidationError {
	var ve validator.ValidationErrors

	if !errors.As(err, &ve) {
		return nil
	}

	result := make([]ValidationError, 0, len(ve))

	for _, fieldErr := range ve {
		result = append(result, ValidationError{
			Property: fieldErr.Field(),
			Tag:      fieldErr.Tag(),
			Value:    fieldErr.Param(),
			Message:  getValidationMessage(fieldErr),
		})
	}

	return &result
}

func getValidationMessage(fieldErr validator.FieldError) string {
	param := fieldErr.Param()
	if param == "" {
		return lang.Trans("validation", fieldErr.Tag())
	}
	return lang.Trans("validation", fieldErr.Tag(), fieldErr.Param())
}
