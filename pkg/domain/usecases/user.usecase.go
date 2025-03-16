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
	Create(CreateRequest) (CreateResponse, *utils.HTTPError)
	GetByID(GetByIDRequest) (GetByIDResponse, *utils.HTTPError)
	GetAll(GetAllRequest) (GetAllResponse, *utils.HTTPError)
	Delete(DeleteRequest) (DeleteResponse, *utils.HTTPError)
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
				fmt.Errorf("[user_uc:GetAccessToken] %w: (%v)", repositories.ErrUserNotFound, err))
		} else {
			e = utils.NewHTTPError(
				utils.StatusInternalServerError,
				"Internal server error",
				"Error during authentication",
				fmt.Errorf("[user_uc:GetAccessToken] %w: (%v)", repositories.ErrDatabase, err))
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
		return GetAccessTokenResponse{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "error during token generation", err)
	}

	return GetAccessTokenResponse{
		Token: accessToken,
	}, nil
}

//
// ======== Create ========
//

type CreateRequest struct {
	Email     vo.Email
	Password  vo.Password
	Lastname  string
	Firstname string
}

type CreateResponse struct {
	entities.User
}

// Create a new user.
func (uc userUseCase) Create(req CreateRequest) (CreateResponse, *utils.HTTPError) {
	// Hash password
	hashedPassword, err := req.Password.HashUserPassword()
	if err != nil {
		return CreateResponse{}, utils.NewHTTPError(utils.StatusInternalServerError, "Error when hashing password", nil, err)
	}
	password, err := vo.NewPassword(hashedPassword)
	if err != nil {
		return CreateResponse{}, utils.NewHTTPError(utils.StatusInternalServerError, "Error when creating password", nil, err)
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
		return CreateResponse{}, utils.NewHTTPError(
			utils.StatusInternalServerError,
			"Internal server error",
			"Error during user creation",
			fmt.Errorf("[user_uc:Create] %w: (%v)", repositories.ErrCreatingUser, err))
	}

	return CreateResponse{
		User: respoRes.User,
	}, nil
}

//
// ======== Get by ID ========
//

// GetByIDRequest is the data transfer object for the GetByID method request.
type GetByIDRequest struct {
	ID entities.UserID
}

// GetByIDResponse is the data transfer object for the GetByID method response.
type GetByIDResponse struct {
	entities.User
}

// GetByID returns a user by its ID.
func (uc userUseCase) GetByID(req GetByIDRequest) (GetByIDResponse, *utils.HTTPError) {
	res, err := uc.userRepository.GetByID(repositories.GetByIDRequest{ID: req.ID})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(
				utils.StatusNotFound,
				"No user found",
				nil,
				fmt.Errorf("[user_uc:GetByID] %w: (%v)", repositories.ErrUserNotFound, err))
		} else {
			e = utils.NewHTTPError(
				utils.StatusInternalServerError,
				"Internal server error",
				"Error when getting user",
				fmt.Errorf("[user_uc:GetByID] %w: (%v)", repositories.ErrGettingUser, err))
		}
		return GetByIDResponse{}, e
	}

	return GetByIDResponse{
		User: res.User,
	}, nil
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
	Data  []entities.User
	Total int
}

// GetAll returns all users (pagination).
func (uc userUseCase) GetAll(req GetAllRequest) (GetAllResponse, *utils.HTTPError) {
	// Get total users
	resTotal, err := uc.userRepository.CountAll(repositories.CountAllRequest{Deleted: req.Deleted})
	if err != nil {
		return GetAllResponse{}, utils.NewHTTPError(
			utils.StatusInternalServerError,
			"Internal server error",
			"Error when getting users",
			fmt.Errorf("[user_uc:GetAll] %w: (%v)", repositories.ErrCountingUsers, err))
	}
	total := resTotal.Total

	users := []entities.User{}
	if total > 0 {
		// Get users
		resUsers, err := uc.userRepository.GetAll(repositories.GetAllRequest{Pagination: req.Pagination, Deleted: req.Deleted})
		if err != nil {
			return GetAllResponse{}, utils.NewHTTPError(
				utils.StatusInternalServerError,
				"Internal server error",
				"Error when getting users",
				fmt.Errorf("[user_uc:GetAll] %w: (%v)", repositories.ErrGettingUsers, err))
		}

		users = resUsers.Users
	}

	return GetAllResponse{
		Data:  users,
		Total: total,
	}, nil
}

//
// ======== Delete ========
//

// DeleteRequest is the data transfer object for the DeleteD method request.
type DeleteRequest struct {
	ID entities.UserID
}

// DeleteResponse is the data transfer object for the DeleteD method response.
type DeleteResponse struct{}

// Delete a user by its ID.
func (uc userUseCase) Delete(req DeleteRequest) (DeleteResponse, *utils.HTTPError) {
	_, err := uc.userRepository.Delete(repositories.DeleteRequest{ID: req.ID})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(
				utils.StatusNotFound,
				"No user found",
				nil,
				fmt.Errorf("[user_uc:GetByID] %w: (%v)", repositories.ErrUserNotFound, err))
		} else {
			e = utils.NewHTTPError(
				utils.StatusInternalServerError,
				"Internal server error",
				"Error when deleting user",
				fmt.Errorf("[user_uc:GetByID] %w: (%v)", repositories.ErrDeletingUser, err))
		}
		return DeleteResponse{}, e
	}

	return DeleteResponse{}, nil
}
