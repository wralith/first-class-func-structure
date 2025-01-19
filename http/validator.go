package http

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	Validator *validator.Validate
}

// Validator needs to implement the Validate method
func (v *StructValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}
