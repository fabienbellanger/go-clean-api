package entities

import (
	vo "go-clean-api/pkg/domain/value_objects"
)

// UserID is a type for user ID
type UserID = vo.ID

// User is a struct that represents a user
type User struct {
	ID        UserID
	Email     vo.Email
	Password  vo.Password
	Lastname  string
	Firstname string
	CreatedAt vo.Time
	UpdatedAt vo.Time
	DeletedAt *vo.Time
}
