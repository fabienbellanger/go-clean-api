package values_objects

import (
	"github.com/google/uuid"
)

type ID struct {
	Value uuid.UUID `validate:"required,uuid"`
}

// String returns the ID value
func (id *ID) String() string {
	return id.Value.String()
}

// NewID creates a new ID
func NewID() ID {
	return ID{Value: uuid.New()}
}

// NewIDFrom creates a new ID from string
func NewIDFrom(value string) (ID, error) {
	uid, err := uuid.Parse(value)
	if err != nil {
		return ID{}, err
	}
	id := ID{Value: uid}

	return id, nil
}
