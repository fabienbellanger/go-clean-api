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

	// ErrCountingUsers is the error returned when counting users.
	ErrCountingUsers = errors.New("error when counting users")

	// ErrGettingUsers is the error returned when getting users.
	ErrGettingUsers = errors.New("error when getting users")

	// ErrGettingUser is the error returned when getting user.
	ErrGettingUser = errors.New("error when getting user")

	// ErrCreatingUser is the error returned when creating user.
	ErrCreatingUser = errors.New("error when creating user")

	// ErrDeletingUser is the error returned when deleting user.
	ErrDeletingUser = errors.New("error when deleting user")
)

// User is the interface that wraps the basic methods to interact with the user repository.
type User interface {
	GetByEmail(GetByEmailRequest) (GetByEmailResponse, error)
	Create(CreateUserRequest) (CreateUserResponse, error)
	GetByID(GetByIDRequest) (GetByIDResponse, error)
	GetAll(GetAllRequest) (GetAllResponse, error)
	CountAll(CountAllRequest) (CountAllResponse, error)
	Delete(DeleteRequest) (DeleteResponse, error)
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

//
// ======== GetAll ========
//

// GetAllRequest is the data transfer object for the GetAll method request.
type GetAllRequest struct {
	Pagination vo.Pagination
	Deleted    bool
}

// GetAllResponse is the data transfer object for the GetAll method response.
type GetAllResponse struct {
	Users []entities.User
}

//
// ======== CountAll ========
//

// CountAllRequest is the data transfer object for the CountAll method request.
type CountAllRequest struct {
	Deleted bool
}

// CountAllResponse is the data transfer object for the CountAll method response.
type CountAllResponse struct {
	Total int
}

//
// ======== Delete ========
//

// DeleteRequest is the data transfer object for the Delete method request.
type DeleteRequest struct {
	ID entities.UserID
}

// DeleteResponse is the data transfer object for the Delete method response.
type DeleteResponse struct{}
