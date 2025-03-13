package values_objects

import (
	"github.com/google/uuid"
)

type ID struct {
	value uuid.UUID `validate:"required,uuid"`
}

// String returns the ID value
func (id *ID) String() string {
	return id.value.String()
}

// Value returns the ID value
func (id *ID) Value() uuid.UUID {
	return id.value
}

// NewID creates a new ID
func NewID() ID {
	return ID{value: uuid.New()}
}

// NewIDFrom creates a new ID from string
func NewIDFrom(value string) (ID, error) {
	uid, err := uuid.Parse(value)
	if err != nil {
		return ID{}, err
	}
	id := ID{value: uid}

	return id, nil
}
