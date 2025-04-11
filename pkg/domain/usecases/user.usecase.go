package usecases

import (
	"errors"
	"fmt"
	"go-clean-api/pkg"
	"go-clean-api/pkg/domain/entities"
	"go-clean-api/pkg/domain/repositories"
	vo "go-clean-api/pkg/domain/value_objects"
	"go-clean-api/utils"
	"time"
)

// User is an interface for user use cases.
type User interface {
	GetAccessToken(GetAccessTokenRequest) (GetAccessTokenResponse, *utils.HTTPError)
	Create(CreateUserRequest) (CreateUserResponse, *utils.HTTPError)
	GetByID(GetUserByIDRequest) (GetUserByIDResponse, *utils.HTTPError)
	GetAll(GetAllUsersRequest) (GetAllUsersResponse, *utils.HTTPError)
	Delete(DeleteRestoreUserRequest) (DeleteRestoreUserResponse, *utils.HTTPError)
	Restore(DeleteRestoreUserRequest) (DeleteRestoreUserResponse, *utils.HTTPError)
}

type userUseCase struct {
	jwtConfig      pkg.ConfigJWT
	userRepository repositories.User
}

// NewUser returns a new User use case
func NewUser(userRepository repositories.User, jwtConfig pkg.ConfigJWT) User {
	return &userUseCase{jwtConfig, userRepository}
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
func (uc userUseCase) GetAccessToken(req GetAccessTokenRequest) (GetAccessTokenResponse, *utils.HTTPError) {
	// Get user ID and password from the email
	userRepo, err := uc.userRepository.GetByEmail(repositories.GetByEmailRequest{Email: req.Email})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(
				utils.StatusUnauthorized,
				"Unauthorizedd",
				nil,
				fmt.Errorf("[user_uc:GetAccessToken] %v", err))
		} else {
			e = utils.NewHTTPError(
				utils.StatusInternalServerError,
				"Internal server error",
				"Error during authentication",
				fmt.Errorf("[user_uc:GetAccessToken] %v", err))
		}
		return GetAccessTokenResponse{}, e
	}

	// Compare the password
	if userRepo.Password.Verify(req.Password.Value()) != nil {
		return GetAccessTokenResponse{}, utils.NewHTTPError(utils.StatusUnauthorized, "Unauthorized", nil, nil)
	}

	// Generate a token
	accessToken, err := entities.NewAccessToken(userRepo.ID, uc.jwtConfig)
	if err != nil {
		return GetAccessTokenResponse{}, utils.NewHTTPError(
			utils.StatusInternalServerError,
			"Internal server error",
			"error during token generation",
			fmt.Errorf("[user_uc:GetAccessToken] token generation error: %v", err))
	}

	return GetAccessTokenResponse{
		Token: accessToken,
	}, nil
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
func (uc userUseCase) Create(req CreateUserRequest) (CreateUserResponse, *utils.HTTPError) {
	// Hash password
	hashedPassword, err := req.Password.HashUserPassword()
	if err != nil {
		return CreateUserResponse{}, utils.NewHTTPError(utils.StatusInternalServerError, "Error when hashing password", nil, err)
	}
	password, err := vo.NewPassword(hashedPassword)
	if err != nil {
		return CreateUserResponse{}, utils.NewHTTPError(utils.StatusInternalServerError, "Error when creating password", nil, err)
	}

	// Add user to the database
	now := vo.NewTime(time.Now(), nil)
	respoRes, err := uc.userRepository.Create(repositories.CreateUserRequest{
		ID:        vo.NewID(),
		Email:     req.Email,
		Password:  password,
		Lastname:  req.Lastname,
		Firstname: req.Firstname,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return CreateUserResponse{}, utils.NewHTTPError(
			utils.StatusInternalServerError,
			"Internal server error",
			"Error during user creation",
			fmt.Errorf("[user_uc:Create] %v", err))
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
func (uc userUseCase) GetByID(req GetUserByIDRequest) (GetUserByIDResponse, *utils.HTTPError) {
	res, err := uc.userRepository.GetByID(repositories.GetByIDRequest{ID: req.ID})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(
				utils.StatusNotFound,
				"No user found",
				nil,
				fmt.Errorf("[user_uc:GetByID] %v", err))
		} else {
			e = utils.NewHTTPError(
				utils.StatusInternalServerError,
				"Internal server error",
				"Error when getting user",
				fmt.Errorf("[user_uc:GetByID] %v", err))
		}
		return GetUserByIDResponse{}, e
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
func (uc userUseCase) GetAll(req GetAllUsersRequest) (GetAllUsersResponse, *utils.HTTPError) {
	// Get total users
	resTotal, err := uc.userRepository.CountAll(repositories.CountAllRequest{Deleted: req.Deleted})
	if err != nil {
		return GetAllUsersResponse{}, utils.NewHTTPError(
			utils.StatusInternalServerError,
			"Internal server error",
			"Error when getting users",
			fmt.Errorf("[user_uc:GetAll] %v", err))
	}
	total := resTotal.Total

	users := []entities.User{}
	if total > 0 {
		// Get users
		resUsers, err := uc.userRepository.GetAll(repositories.GetAllRequest{Pagination: req.Pagination, Deleted: req.Deleted})
		if err != nil {
			return GetAllUsersResponse{}, utils.NewHTTPError(
				utils.StatusInternalServerError,
				"Internal server error",
				"Error when getting users",
				fmt.Errorf("[user_uc:GetAll] %v", err))
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
func (uc userUseCase) Delete(req DeleteRestoreUserRequest) (DeleteRestoreUserResponse, *utils.HTTPError) {
	_, err := uc.userRepository.Delete(repositories.DeleteRestoreRequest{ID: req.ID})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(
				utils.StatusNotFound,
				"No user found",
				nil,
				fmt.Errorf("[user_uc:Delete] %v", err))
		} else {
			e = utils.NewHTTPError(
				utils.StatusInternalServerError,
				"Internal server error",
				"Error when deleting user",
				fmt.Errorf("[user_uc:Delete] %v", err))
		}
		return DeleteRestoreUserResponse{}, e
	}

	return DeleteRestoreUserResponse{}, nil
}

// Restore a user by its ID.
func (uc userUseCase) Restore(req DeleteRestoreUserRequest) (DeleteRestoreUserResponse, *utils.HTTPError) {
	_, err := uc.userRepository.Restore(repositories.DeleteRestoreRequest{ID: req.ID})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(
				utils.StatusNotFound,
				"No user found",
				nil,
				fmt.Errorf("[user_uc:Restore] %v", err))
		} else {
			e = utils.NewHTTPError(
				utils.StatusInternalServerError,
				"Internal server error",
				"Error when restoring user",
				fmt.Errorf("[user_uc:Restore] %v", err))
		}
		return DeleteRestoreUserResponse{}, e
	}

	return DeleteRestoreUserResponse{}, nil
}
