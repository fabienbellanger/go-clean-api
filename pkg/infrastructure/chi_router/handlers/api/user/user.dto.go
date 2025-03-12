package user

import (
	"go-clean-api/pkg/domain/usecases"
	vo "go-clean-api/pkg/domain/value_objects"
	"time"
)

type GetAccessTokenRequest struct {
	Email    string `json:"email" xml:"email" form:"email"`
	Password string `json:"password" xml:"password" form:"password"`
}

func (r GetAccessTokenRequest) ToUseCase() (usecases.GetAccessTokenRequest, error) {
	email, err := vo.NewEmail(r.Email)
	if err != nil {
		return usecases.GetAccessTokenRequest{}, err
	}

	password, err := vo.NewPassword(r.Password)
	if err != nil {
		return usecases.GetAccessTokenRequest{}, err
	}

	return usecases.GetAccessTokenRequest{
		Email:    email,
		Password: password,
	}, nil
}

type GetAccessTokenResponse struct {
	AccessToken          string    `json:"access_token" xml:"access_token"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at" xml:"access_token_expired_at"`
}
