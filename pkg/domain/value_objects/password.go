package values_objects

import (
	"go-clean-api/utils"

	"golang.org/x/crypto/bcrypt"
)

// Password represents an password value object
type Password struct {
	value string
}

func (p *Password) String() string {
	return p.Value()
}

// Value returns the password value
func (p *Password) Value() string {
	return p.value
}

// NewPassword creates a new password
func NewPassword(value string) (Password, error) {
	p := Password{value: value}

	err := p.Validate()
	if err != nil {
		return Password{}, &err
	}

	return p, nil
}

// Validate checks if a struct is valid and returns an array of errors
func (p *Password) Validate() utils.ValidatorErrors {
	return utils.ValidateVar(p.value, "password", "required,min=8")
}

// HashUserPassword hashes a password
func (p *Password) HashUserPassword() (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	return string(passwordBytes), err
}

// Verify checks if the password is correct
func (p *Password) Verify(hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.value), []byte(hashedPassword))
}
