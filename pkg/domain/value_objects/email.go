package values_objects

import "go-clean-api/utils"

// Email represents an email value object
type Email struct {
	Value string `validate:"required,email"`
}

// String returns the email value
func (e *Email) String() string {
	return e.Value
}

// NewEmail creates a new email
func NewEmail(value string) (Email, error) {
	e := Email{Value: value}

	err := e.Validate()
	if err != nil {
		return Email{}, &err
	}

	return e, nil
}

// Validate checks if a struct is valid and returns an array of errors
func (e *Email) Validate() utils.ValidatorErrors {
	return utils.ValidateStruct(e)
}
