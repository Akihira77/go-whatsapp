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
	errs := v.validate.Struct(data)
	if errs == nil {
		return nil
	}

	vErrs := make(map[string]string, 0)
	for _, v := range errs.(validator.ValidationErrors) {
		var e error
		switch v.Tag() {
		case "required":
			e = fmt.Errorf("Field '%s' cannot be empty", v.Field())
		case "email":
			e = fmt.Errorf("Field '%s' must be a valid email address", v.Field())
		case "eth_addr":
			e = fmt.Errorf("Field '%s' must  be a valid Ethereum address", v.Field())
		case "len":
			e = fmt.Errorf("Field '%s' must be exactly %v characters long", v.Field(), v.Param())
		default:
			e = fmt.Errorf("Field '%s' must satisfy '%s' '%v' criteria", v.Field(), v.Tag(), v.Param())
		}

		vErrs[strings.ToLower(v.Field())] = e.Error()
	}

	return vErrs
}
