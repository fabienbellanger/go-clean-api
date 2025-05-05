package models

import (
	"fmt"
	"go-clean-api/pkg/domain/entities"
	"go-clean-api/pkg/domain/repositories"
	vo "go-clean-api/pkg/domain/value_objects"
)

// GetUserByEmail is the data transfer object for the GetUserByEmail method request.
type GetUserByEmail struct {
	ID       string `db:"id"`
	Password string `db:"password"`
}

// ToRepository converts the model to repository response
// TODO: Add tests
func (u GetUserByEmail) Repository() (repositories.GetByEmailResponse, error) {
	id, err := vo.NewIDFrom(u.ID)
	if err != nil {
		return repositories.GetByEmailResponse{}, fmt.Errorf("[models:GetUserByEmail %w: %s]", ErrIDFromString, err)
	}

	password, err := vo.NewPassword(u.Password)
	if err != nil {
		return repositories.GetByEmailResponse{}, fmt.Errorf("[models:GetUserByEmail %w: %s]", ErrPasswordFromString, err)
	}
	return repositories.GetByEmailResponse{
		ID:       id,
		Password: password,
	}, nil
}

// User is the data transfer object for the User entity
type User struct {
	ID        string  `db:"id"`
	Email     string  `db:"email"`
	Lastname  string  `db:"lastname"`
	Firstname string  `db:"firstname"`
	CreatedAt string  `db:"created_at"` // Format YYYY-MM-DD HH:MM:SS
	UpdatedAt string  `db:"updated_at"` // Format YYYY-MM-DD HH:MM:SS
	DeletedAt *string `db:"deleted_at"` // Format YYYY-MM-DD HH:MM:SS
}

// Entity converts the user model to entity
// TODO: Add tests
func (u User) Entity() (user entities.User, err error) {
	id, errID := vo.NewIDFrom(u.ID)
	if errID != nil {
		err = fmt.Errorf("[models:User:Entity %w: %s]", ErrIDFromString, errID)
		return
	}

	email, errEmail := vo.NewEmail(u.Email)
	if errEmail != nil {
		err = fmt.Errorf("[models:User:Entity %w: %s]", ErrEmailFromString, errEmail)
		return
	}

	createdAt, errDateTime := vo.ParseRFC3339(u.CreatedAt, nil)
	if errDateTime != nil {
		err = fmt.Errorf("[models:User:Entity %w: %s]", errDateTime, errDateTime)
		return
	}

	updatedAt, errDateTime := vo.ParseRFC3339(u.UpdatedAt, nil)
	if errDateTime != nil {
		err = fmt.Errorf("[models:User:Entity %w: %s]", errDateTime, errDateTime)
		return
	}

	var deletedAt *vo.Time
	if u.DeletedAt != nil {
		d, errDateTime := vo.ParseRFC3339(*u.DeletedAt, nil)
		if errDateTime != nil {
			err = fmt.Errorf("[models:User:Entity %w: %s]", errDateTime, errDateTime)
			return
		}
		deletedAt = &d
	}

	user = entities.User{
		ID:        id,
		Email:     email,
		Lastname:  u.Lastname,
		Firstname: u.Firstname,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}

	return
}
