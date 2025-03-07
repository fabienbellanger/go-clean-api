package values_objects

import (
	"go-clean-api/utils"

	"golang.org/x/crypto/bcrypt"
)

// Password represents an password value object
type Password struct {
	Value string `validate:"required,min=8"`
}

// String returns the password value
func (p *Password) String() string {
	return p.Value
}

// NewPassword creates a new password
func NewPassword(value string) (Password, error) {
	p := Password{Value: value}

	err := p.Validate()
	if err != nil {
		return Password{}, &err
	}

	return p, nil
}

// Validate checks if a struct is valid and returns an array of errors
func (p *Password) Validate() utils.ValidatorErrors {
	return utils.ValidateStruct(p)
}

// HashUserPassword hashes a password
func (p *Password) HashUserPassword() (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(p.Value), bcrypt.DefaultCost)
	return string(passwordBytes), err
}

// Verify checks if the password is correct
func (p *Password) Verify(hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.Value), []byte(hashedPassword))
}
