package usecases

import (
	"errors"
	"go-clean-api/pkg"
	"go-clean-api/pkg/domain/entities"
	"go-clean-api/pkg/domain/repositories"
	vo "go-clean-api/pkg/domain/value_objects"
	"go-clean-api/utils"
)

// User is an interface for user use cases.
type User interface {
	GetAccessToken(GetAccessTokenRequest) (GetAccessTokenResponse, *utils.HTTPError)
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
	if userRepo.Password.Verify(req.Password.Value) != nil {
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
