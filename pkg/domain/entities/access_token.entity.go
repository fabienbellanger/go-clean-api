package entities

import (
	vo "go-clean-api/pkg/domain/value_objects"
)

// AccessToken is a struct that represents a JWT access token
type AccessToken struct {
	Token     string
	ExpiredAt vo.Time
}
