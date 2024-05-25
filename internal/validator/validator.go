package validator

import "github.com/go-playground/validator/v10"

func NewCustomValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate
}
