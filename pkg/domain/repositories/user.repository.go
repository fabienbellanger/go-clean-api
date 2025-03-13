package repositories

import (
	"errors"
	"go-clean-api/pkg/domain/entities"
	vo "go-clean-api/pkg/domain/value_objects"
)

var (
	// ErrUserNotFound is the error returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrDatabase is the error returned when a database error occurs.
	ErrDatabase = errors.New("database error")

	// ErrConvertFromModel is the error returned when a model cannot be converted to repository response.
	ErrConvertFromModel = errors.New("error converting model to repository response")
)

// User is the interface that wraps the basic methods to interact with the user repository.
type User interface {
	GetByEmail(GetByEmailRequest) (GetByEmailResponse, error)
	Create(CreateUserRequest) (CreateUserResponse, error)
	GetByID(GetByIDRequest) (GetByIDResponse, error)
}

//
// ======== GetByEmail ========
//

// GetByEmailRequest is the data transfer object for the GetByEmail method request.
type GetByEmailRequest struct {
	Email vo.Email
}

// GetByEmailResponse is the data transfer object for the GetByEmail method response.
type GetByEmailResponse struct {
	ID       entities.UserID
	Password vo.Password
}

//
// ======== Create ========
//

type CreateUserRequest struct {
	ID        entities.UserID
	Email     vo.Email
	Password  vo.Password
	Lastname  string
	Firstname string
	CreatedAt vo.Time
	UpdatedAt vo.Time
}

type CreateUserResponse struct {
	entities.User
}

//
// ======== GetByID ========
//

// GetByIDRequest is the data transfer object for the GetByID method request.
type GetByIDRequest struct {
	ID entities.UserID
}

// GetByIDResponse is the data transfer object for the GetByID method response.
type GetByIDResponse struct {
	entities.User
}
