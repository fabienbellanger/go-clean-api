package values_objects

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewIDFrom(t *testing.T) {
	type result struct {
		id     ID
		is_err bool
	}

	tests := []struct {
		value  string
		wanted result
	}{
		{
			value: "550e8400-e29b-41d4-a716-446655440000",
			wanted: result{
				id:     ID{Value: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
				is_err: false,
			},
		},
		{
			value: "f5ds415f4",
			wanted: result{
				id:     ID{},
				is_err: true,
			},
		},
		{
			value: "",
			wanted: result{
				id:     ID{},
				is_err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, err := NewIDFrom(tt.value)
			log.Printf("Error: %v\n", err)
			if err != nil {
				assert.Equal(t, tt.wanted.is_err, true)
			} else {
				assert.Equal(t, got, tt.wanted.id)
			}
		})
	}
}
