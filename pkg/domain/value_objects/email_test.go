package values_objects

import (
	"testing"

	"go-clean-api/utils"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	type result struct {
		email Email
		err   error
	}

	var e1 utils.ValidatorErrors
	e1 = append(e1, utils.ValidatorError{
		Field: "email",
		Tag:   "email",
		Value: "",
	})
	var e2 utils.ValidatorErrors
	e2 = append(e2, utils.ValidatorError{
		Field: "email",
		Tag:   "required",
		Value: "",
	})

	tests := []struct {
		value  string
		wanted result
	}{
		{
			value: "toto@gmail.com",
			wanted: result{
				email: Email{value: "toto@gmail.com"},
				err:   nil,
			},
		},
		{
			value: "bad",
			wanted: result{
				email: Email{},
				err:   &e1,
			},
		},
		{
			value: "",
			wanted: result{
				email: Email{},
				err:   &e2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, err := NewEmail(tt.value)

			if err != nil {
				assert.Equal(t, err, tt.wanted.err)
			} else {
				assert.Equal(t, got, tt.wanted.email)
			}
		})
	}
}
