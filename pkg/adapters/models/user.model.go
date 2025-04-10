package models

import (
	"go-clean-api/pkg/domain/entities"
	"go-clean-api/pkg/domain/repositories"
	vo "go-clean-api/pkg/domain/value_objects"
)

// GetByEmail is the data transfer object for the GetByEmail method request.
type GetByEmail struct {
	ID       string `db:"id"`
	Password string `db:"password"`
}

// ToRepository converts the model to repository response
// TODO: Add tests
func (u GetByEmail) Repository() (repositories.GetByEmailResponse, error) {
	id, err := vo.NewIDFrom(u.ID)
	if err != nil {
		return repositories.GetByEmailResponse{}, err
	}

	password, err := vo.NewPassword(u.Password)
	if err != nil {
		return repositories.GetByEmailResponse{}, err
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
	id, err := vo.NewIDFrom(u.ID)
	if err != nil {
		return
	}

	email, err := vo.NewEmail(u.Email)
	if err != nil {
		return
	}

	createdAt, err := vo.ParseRFC3339(u.CreatedAt, nil)
	if err != nil {
		return
	}

	updatedAt, err := vo.ParseRFC3339(u.UpdatedAt, nil)
	if err != nil {
		return
	}

	var deletedAt *vo.Time
	if u.DeletedAt != nil {
		d, e := vo.ParseRFC3339(*u.DeletedAt, nil)
		if e != nil {
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

	return user, nil
}
