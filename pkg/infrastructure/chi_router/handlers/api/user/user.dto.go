package user

import (
	"go-clean-api/pkg/domain/usecases"
	vo "go-clean-api/pkg/domain/value_objects"
	"time"
)

type GetAccessTokenRequest struct {
	Email string `json:"email" xml:"email" form:"email"`
	Passw string `json:"password" xml:"password" form:"password"`
}

func (r GetAccessTokenRequest) ToUseCase() (usecases.GetAccessTokenRequest, error) {
	email, err := vo.NewEmail(r.Email)
	if err != nil {
		return usecases.GetAccessTokenRequest{}, err
	}

	password, err := vo.NewPassword(r.Passw)
	if err != nil {
		return usecases.GetAccessTokenRequest{}, err
	}

	return usecases.GetAccessTokenRequest{
		Email:    email,
		Password: password,
	}, nil
}

type GetAccessTokenResponse struct {
	AccessToken string    `json:"access_token" xml:"access_token"`
	ExpireAt    time.Time `json:"expires_at" xml:"expires_at"`
}
