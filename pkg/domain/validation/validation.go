package validation

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

// ValidatorError represents error validation struct.
type ValidatorError struct {
	Field string `json:"field" xml:"field"`
	Tag   string `json:"tag" xml:"tag"`
	Value string `json:"value" xml:"value"`
}

// ValidatorErrors is a slice of ValidatorError.
type ValidatorErrors []ValidatorError

func (ve *ValidatorErrors) Error() string {
	b, err := json.Marshal(ve)
	if err != nil {
		return "error when marshalling validation errors"
	}
	return string(b)
}

// ValidateStruct checks if a struct is valid and returns an array of errors
// if it is not valid.
func ValidateStruct(s any) (errors ValidatorErrors) {
	validate := validator.New()
	errs := validate.Struct(s)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			errors = append(errors, ValidatorError{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: err.Param(),
			})
		}
	}
	return
}

// ValidateVar checks if a variable is valid and returns an array of errors
// if it is not valid.
func ValidateVar(v any, field, tag string) (errors ValidatorErrors) {
	validate := validator.New()
	errs := validate.Var(v, tag)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			errors = append(errors, ValidatorError{
				Field: field,
				Tag:   err.Tag(),
				Value: err.Param(),
			})
		}
	}
	return
}
