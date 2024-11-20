package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func New() *validator.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation("notblank", validators.NotBlank)
	return validate
}

func Validate(s interface{}) error {
	return New().Struct(s)
}
