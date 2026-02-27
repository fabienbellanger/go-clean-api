package usecases

import (
	"errors"
	"fmt"
	domainerr "go-clean-api/pkg/domain/errors"
	"go-clean-api/pkg/domain/entities"
	"go-clean-api/pkg/domain/repositories"
	"go-clean-api/pkg/domain/services"
	vo "go-clean-api/pkg/domain/value_objects"
	"time"
)

var (
	ErrInvalidPassword     = errors.New("invalid password")
	ErrHashPassword        = errors.New("error when hashing password")
	ErrAccessTokenCreation = errors.New("error when creating access token")
	ErrUserCreation        = errors.New("error when creating user")
)

// User is an interface for user use cases.
type User interface {
	GetAccessToken(GetAccessTokenRequest) (GetAccessTokenResponse, error)
	Create(CreateUserRequest) (CreateUserResponse, error)
	GetByID(GetUserByIDRequest) (GetUserByIDResponse, error)
	GetAll(GetAllUsersRequest) (GetAllUsersResponse, error)
	Delete(DeleteRestoreUserRequest) (DeleteRestoreUserResponse, error)
	Restore(DeleteRestoreUserRequest) (DeleteRestoreUserResponse, error)
}

type userUseCase struct {
	tokenGenerator services.TokenGenerator
	userRepository repositories.User
}

// NewUser returns a new User use case
func NewUser(userRepository repositories.User, tokenGenerator services.TokenGenerator) User {
	return &userUseCase{tokenGenerator, userRepository}
}

//
// ======== GetAccessToken ========
//

// GetAccessTokenRequest is the data transfer object for the GetAccessToken method request.
type GetAccessTokenRequest struct {
	Email    vo.Email
	Password vo.Password
}

// GetAccessTokenResponse is the data transfer object for the GetAccessToken method response.
type GetAccessTokenResponse struct {
	Token entities.AccessToken
}

// GetAccessToken returns an access token from user email and password.
func (uc userUseCase) GetAccessToken(req GetAccessTokenRequest) (res GetAccessTokenResponse, err error) {
	// Get user ID and password from the email
	userRepo, errRepo := uc.userRepository.GetByEmail(repositories.GetByEmailRequest{Email: req.Email})
	if errRepo != nil {
		if errors.Is(errRepo, domainerr.ErrNotFound) {
			err = fmt.Errorf("[user_uc:GetAccessToken %w: %s]", domainerr.ErrNotFound, errRepo)
		} else {
			err = fmt.Errorf("[user_uc:GetAccessToken %w: %s]", domainerr.ErrDatabase, errRepo)
		}
		return
	}

	// Compare the password
	if userRepo.Password.Verify(req.Password.Value()) != nil {
		err = fmt.Errorf("[user_uc:GetAccessToken %w]", ErrInvalidPassword)
		return
	}

	// Generate a token
	accessToken, errToken := uc.tokenGenerator.Generate(userRepo.ID)
	if errToken != nil {
		err = fmt.Errorf("[user_uc:GetAccessToken %w: %s]", ErrAccessTokenCreation, errToken)
		return
	}

	res.Token = accessToken

	return
}

//
// ======== Create ========
//

type CreateUserRequest struct {
	Email     vo.Email
	Password  vo.Password
	Lastname  string
	Firstname string
}

type CreateUserResponse struct {
	entities.User
}

// Create a new user.
func (uc userUseCase) Create(req CreateUserRequest) (res CreateUserResponse, err error) {
	// Hash password
	hashedPassword, errHash := req.Password.HashUserPassword()
	if errHash != nil {
		err = fmt.Errorf("[user_uc:Create %w: %s]", ErrHashPassword, errHash)
		return
	}
	password, errPassword := vo.NewPassword(hashedPassword)
	if errPassword != nil {
		err = fmt.Errorf("[user_uc:Create %w: %s]", ErrInvalidPassword, errPassword)
		return
	}

	// Add user to the database
	now := vo.NewTime(time.Now(), nil)
	respoRes, errRepo := uc.userRepository.Create(repositories.CreateUserRequest{
		ID:        vo.NewID(),
		Email:     req.Email,
		Password:  password,
		Lastname:  req.Lastname,
		Firstname: req.Firstname,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if errRepo != nil {
		err = fmt.Errorf("[user_uc:Create %w: %s]", ErrUserCreation, errRepo)
		return
	}

	return CreateUserResponse{
		User: respoRes.User,
	}, nil
}

//
// ======== Get by ID ========
//

// GetUserByIDRequest is the data transfer object for the GetByID method request.
type GetUserByIDRequest struct {
	ID entities.UserID
}

// GetUserByIDResponse is the data transfer object for the GetByID method response.
type GetUserByIDResponse struct {
	entities.User
}

// GetByID returns a user by its ID.
func (uc userUseCase) GetByID(req GetUserByIDRequest) (GetUserByIDResponse, error) {
	res, err := uc.userRepository.GetByID(repositories.GetByIDRequest{ID: req.ID})
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			return GetUserByIDResponse{}, fmt.Errorf("[user_uc:GetByID %w: %s]", domainerr.ErrNotFound, err)
		}
		return GetUserByIDResponse{}, fmt.Errorf("[user_uc:GetByID %w: %s]", domainerr.ErrDatabase, err)
	}

	return GetUserByIDResponse{
		User: res.User,
	}, nil
}

//
// ======== GetAll ========
//

// GetAllUsersRequest is the data transfer object for the GetAll method request.
type GetAllUsersRequest struct {
	Pagination vo.Pagination
	Deleted    bool
}

// GetAllUsersResponse is the data transfer object for the GetAll method response.
type GetAllUsersResponse struct {
	Data  []entities.User
	Total int64
}

// GetAll returns all users (pagination).
func (uc userUseCase) GetAll(req GetAllUsersRequest) (res GetAllUsersResponse, err error) {
	// Get total users
	resTotal, errTotal := uc.userRepository.CountAll(repositories.CountAllRequest{Deleted: req.Deleted})
	if errTotal != nil {
		err = fmt.Errorf("[user_uc:GetAll %w: %s]", domainerr.ErrDatabase, errTotal)
		return
	}
	total := resTotal.Total

	users := []entities.User{}
	if total > 0 {
		// Get users
		resUsers, errUsers := uc.userRepository.GetAll(repositories.GetAllRequest{Pagination: req.Pagination, Deleted: req.Deleted})
		if errUsers != nil {
			err = fmt.Errorf("[user_uc:GetAll %w: %s]", domainerr.ErrDatabase, errUsers)
			return
		}

		users = resUsers.Users
	}

	return GetAllUsersResponse{
		Data:  users,
		Total: total,
	}, nil
}

//
// ======== Delete / Restore ========
//

// DeleteRestoreUserRequest is the data transfer object for the DeleteD method request.
type DeleteRestoreUserRequest struct {
	ID entities.UserID
}

// DeleteRestoreUserResponse is the data transfer object for the DeleteD method response.
type DeleteRestoreUserResponse struct{}

// Delete a user by its ID.
func (uc userUseCase) Delete(req DeleteRestoreUserRequest) (DeleteRestoreUserResponse, error) {
	_, err := uc.userRepository.Delete(repositories.DeleteRestoreRequest{ID: req.ID})
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			return DeleteRestoreUserResponse{}, fmt.Errorf("[user_uc:Delete %w: %s]", domainerr.ErrNotFound, err)
		}
		return DeleteRestoreUserResponse{}, fmt.Errorf("[user_uc:Delete %w: %s]", domainerr.ErrDatabase, err)
	}

	return DeleteRestoreUserResponse{}, nil
}

// Restore a user by its ID.
func (uc userUseCase) Restore(req DeleteRestoreUserRequest) (DeleteRestoreUserResponse, error) {
	_, err := uc.userRepository.Restore(repositories.DeleteRestoreRequest{ID: req.ID})
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			return DeleteRestoreUserResponse{}, fmt.Errorf("[user_uc:Restore %w: %s]", domainerr.ErrNotFound, err)
		}
		return DeleteRestoreUserResponse{}, fmt.Errorf("[user_uc:Restore %w: %s]", domainerr.ErrDatabase, err)
	}

	return DeleteRestoreUserResponse{}, nil
}
