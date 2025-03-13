package user

import (
	"go-clean-api/pkg/domain/usecases"
	vo "go-clean-api/pkg/domain/value_objects"
	"time"
)

//
// ======== GetAccessToken ========
//

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

//
// ======== Create ========
//

type CreateRequest struct {
	Email     string `json:"email" xml:"email" form:"email"`
	Password  string `json:"password" xml:"password" form:"password"`
	Lastname  string `json:"lastname" xml:"lastname" form:"lastname"`
	Firstname string `json:"firstname" xml:"firstname" form:"firstname"`
}

func (r CreateRequest) ToUseCase() (usecases.CreateRequest, error) {
	email, err := vo.NewEmail(r.Email)
	if err != nil {
		return usecases.CreateRequest{}, err
	}

	password, err := vo.NewPassword(r.Password)
	if err != nil {
		return usecases.CreateRequest{}, err
	}

	return usecases.CreateRequest{
		Email:     email,
		Password:  password,
		Lastname:  r.Lastname,
		Firstname: r.Firstname,
	}, nil
}

type CreateResponse struct {
	ID        string `json:"id" xml:"id"`
	Email     string `json:"email" xml:"email"`
	Lastname  string `json:"lastname" xml:"lastname"`
	Firstname string `json:"firstname" xml:"firstname"`
	CreatedAt string `json:"created_at" xml:"created_at"`
	UpdatedAt string `json:"updated_at" xml:"updated_at"`
}
