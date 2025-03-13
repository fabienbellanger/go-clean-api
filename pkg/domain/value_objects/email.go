package values_objects

import (
	"go-clean-api/utils"
)

// Email represents an email value object
type Email struct {
	value string `validate:"required,email"`
}

func (e *Email) String() string {
	return e.Value()
}

// Value returns the email value
func (e *Email) Value() string {
	return e.value
}

// NewEmail creates a new email
func NewEmail(value string) (Email, error) {
	e := Email{value: value}

	err := e.Validate()
	if err != nil {
		return Email{}, &err
	}

	return e, nil
}

// Validate checks if a struct is valid and returns an array of errors
func (e *Email) Validate() utils.ValidatorErrors {
	return utils.ValidateVar(e.value, "email", "required,email")
}
