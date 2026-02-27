package values_objects

import (
	"log"
	"testing"

	"go-clean-api/pkg/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestNewPassword(t *testing.T) {
	type result struct {
		password Password
		err      error
	}

	var e1 validation.ValidatorErrors
	e1 = append(e1, validation.ValidatorError{
		Field: "password",
		Tag:   "min",
		Value: "8",
	})
	var e2 validation.ValidatorErrors
	e2 = append(e2, validation.ValidatorError{
		Field: "password",
		Tag:   "required",
		Value: "",
	})

	tests := []struct {
		value  string
		wanted result
	}{
		{
			value: "password",
			wanted: result{
				password: Password{value: "password"},
				err:      nil,
			},
		},
		{
			value: "bad",
			wanted: result{
				password: Password{},
				err:      &e1,
			},
		},
		{
			value: "",
			wanted: result{
				password: Password{},
				err:      &e2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, err := NewPassword(tt.value)
			log.Printf("Error: %v\n", err)
			if err != nil {
				assert.Equal(t, err, tt.wanted.err)
			} else {
				assert.Equal(t, got, tt.wanted.password)
			}
		})
	}
}
