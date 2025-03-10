package models

import (
	"go-clean-api/pkg/domain/repositories"
	vo "go-clean-api/pkg/domain/value_objects"
)

// GetByEmail is the data transfer object for the GetByEmail method request.
type GetByEmail struct {
	ID       string
	Password string
}

// ToRepository converts the model to repository response
func (u GetByEmail) ToRepository() (repositories.GetByEmailResponse, error) {
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
