package services

import "go-clean-api/pkg/domain/entities"

// TokenGenerator defines the interface for generating access tokens.
type TokenGenerator interface {
	Generate(userID entities.UserID) (entities.AccessToken, error)
}
