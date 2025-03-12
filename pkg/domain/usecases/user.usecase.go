package usecases

import (
	"errors"
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
			e = utils.NewHTTPError(utils.StatusUnauthorized, "Unauthorized", nil, nil)
		} else {
			e = utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "error during authentication", err)
		}
		return GetAccessTokenResponse{}, e
	}

	// Compare the password
	if userRepo.Password.Verify(req.Password.String()) != nil {
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
		return CreateResponse{}, utils.NewHTTPError(utils.StatusInternalServerError, "Error when hashing password", err, nil)
	}
	password, err := vo.NewPassword(hashedPassword)
	if err != nil {
		return CreateResponse{}, utils.NewHTTPError(utils.StatusInternalServerError, "Error when creating password", err, nil)
	}

	// Add user to the database
	now := time.Now()
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
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrDatabase) {
			e = utils.NewHTTPError(utils.StatusInternalServerError, "Error when creating user", nil, nil)
		} else {
			e = utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "error during authentication", err)
		}
		return CreateResponse{}, e
	}

	return CreateResponse{
		User: respoRes.User,
	}, nil
}
