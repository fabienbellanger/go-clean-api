package user

import (
	"go-clean-api/pkg/domain/usecases"
	vo "go-clean-api/pkg/domain/value_objects"
)

type UserResponse struct {
	ID        string `json:"id" xml:"id"`
	Email     string `json:"email" xml:"email"`
	Lastname  string `json:"lastname" xml:"lastname"`
	Firstname string `json:"firstname" xml:"firstname"`
	CreatedAt string `json:"created_at" xml:"created_at"`
	UpdatedAt string `json:"updated_at" xml:"updated_at"`
	DeletedAt string `json:"deleted_at,omitempty" xml:"deleted_at,omitempty"`
}

//
// ======== GetAccessToken ========
//

type GetAccessTokenRequest struct {
	Email    string `json:"email" xml:"email" form:"email"`
	Password string `json:"password" xml:"password" form:"password"`
}

// TODO: Add tests
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
	AccessToken          string `json:"access_token" xml:"access_token"`
	AccessTokenExpiredAt string `json:"access_token_expired_at" xml:"access_token_expired_at"`
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

// TODO: Add tests
func (r CreateRequest) ToUseCase() (usecases.CreateUserRequest, error) {
	email, err := vo.NewEmail(r.Email)
	if err != nil {
		return usecases.CreateUserRequest{}, err
	}

	password, err := vo.NewPassword(r.Password)
	if err != nil {
		return usecases.CreateUserRequest{}, err
	}

	return usecases.CreateUserRequest{
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

//
// ======== Get by ID ========
//

type GetByIDRequest struct {
	ID string
}

func (r GetByIDRequest) ToUseCase() (usecases.GetUserByIDRequest, error) {
	id, err := vo.NewIDFrom(r.ID)
	if err != nil {
		return usecases.GetUserByIDRequest{}, err
	}

	return usecases.GetUserByIDRequest{
		ID: id,
	}, nil
}

type GetByIDResponse struct {
	UserResponse
}

// TODO: Add tests
func (r GetByIDResponse) FromEntity(res usecases.GetUserByIDResponse) GetByIDResponse {
	deletedAt := ""
	if res.DeletedAt != nil {
		deletedAt = res.DeletedAt.RFC3339()
	}

	r.ID = res.ID.String()
	r.Email = res.Email.Value()
	r.Lastname = res.Lastname
	r.Firstname = res.Firstname
	r.CreatedAt = res.CreatedAt.RFC3339()
	r.UpdatedAt = res.UpdatedAt.RFC3339()
	r.DeletedAt = deletedAt

	return r
}

//
// ======== Get all ========
//

type GetAllResponse struct {
	Data  []UserResponse `json:"data" xml:"data"`
	Page  int            `json:"page" xml:"page"`
	Size  int            `json:"size" xml:"size"`
	Total int            `json:"total" xml:"total"`
}

// TODO: Add tests
func (r GetAllResponse) FromEntity(res usecases.GetAllUsersResponse, pagination vo.Pagination) GetAllResponse {
	r.Data = make([]UserResponse, len(res.Data))
	for i, user := range res.Data {
		deletedAt := ""
		if user.DeletedAt != nil {
			deletedAt = user.DeletedAt.RFC3339()
		}

		r.Data[i] = UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email.Value(),
			Lastname:  user.Lastname,
			Firstname: user.Firstname,
			CreatedAt: user.CreatedAt.RFC3339(),
			UpdatedAt: user.UpdatedAt.RFC3339(),
			DeletedAt: deletedAt,
		}
	}

	r.Total = res.Total
	r.Page = pagination.Page()
	r.Size = pagination.Size()

	return r
}

//
// ======== Delete / Restore ========
//

type DeleteRestoreRequest struct {
	ID string
}

func (r DeleteRestoreRequest) ToUseCase() (usecases.DeleteRestoreUserRequest, error) {
	id, err := vo.NewIDFrom(r.ID)
	if err != nil {
		return usecases.DeleteRestoreUserRequest{}, err
	}

	return usecases.DeleteRestoreUserRequest{
		ID: id,
	}, nil
}
