package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidatorError represents error validation struct.
type ValidatorError struct {
	FailedField string
	Tag         string
	Value       string
}

// TODO: Change
func (ve *ValidatorError) Error() string {
	return fmt.Sprintf("Field: %s, Tag: %s, Value: %s", strings.ToLower(ve.FailedField), ve.Tag, ve.Value)
}

// ValidatorErrors is a slice of ValidatorError.
type ValidatorErrors []ValidatorError

// TODO: Change
func (ve *ValidatorErrors) Error() string {
	e := ""
	for i, v := range *ve {
		if i != 0 {
			e += "\n"
		}
		e += v.Error()
	}
	return e
}

// ValidateStruct checks if a struct is valid and returns an array of errors
// if it is not valid.
func ValidateStruct(s any) (errors ValidatorErrors) {
	validate := validator.New()
	errs := validate.Struct(s)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			errors = append(errors, ValidatorError{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Param(),
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
			log.Printf("%#v\n", err)
			errors = append(errors, ValidatorError{
				FailedField: field,
				Tag:         err.Tag(),
				Value:       err.Param(),
			})
		}
	}
	return
}
