package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MyValidator struct {
	validate *validator.Validate
}

func NewMyValidator() *MyValidator {
	return &MyValidator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *MyValidator) Validate(data any) map[string]string {
	err := v.validate.Struct(data)
	if err == nil {
		return nil
	}

	validationErrs := make(map[string]string, 0)
	for _, v := range err.(validator.ValidationErrors) {
		var e error
		switch v.Tag() {
		case "required":
			e = fmt.Errorf("Field '%s' cannot be empty", v.Field())
		case "email":
			e = fmt.Errorf("Field '%s' must be a valid email address", v.Field())
		case "len":
			e = fmt.Errorf("Field '%s' must be exactly %v characters long", v.Field(), v.Param())
		default:
			e = fmt.Errorf("Field '%s' must satisfy '%s' '%v' criteria", v.Field(), v.Tag(), v.Param())
		}

		validationErrs[strings.ToLower(v.Field())] = e.Error()
	}

	return validationErrs
}
