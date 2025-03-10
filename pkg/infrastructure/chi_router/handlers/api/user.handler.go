package api

import (
	"encoding/json"
	"go-clean-api/pkg/domain/usecases"
	"go-clean-api/pkg/infrastructure/chi_router/handlers"
	"go-clean-api/pkg/infrastructure/logger"
	"go-clean-api/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// User handler
type User struct {
	router      chi.Router
	userUseCase usecases.User
	logger      logger.CustomLogger
}

// NewUser returns a new Handler
func NewUser(r chi.Router, l logger.CustomLogger, userUseCase usecases.User) User {
	return User{
		router:      r,
		userUseCase: userUseCase,
		logger:      l,
	}
}

// UserPublicRoutes adds users public routes
func (u *User) UserPublicRoutes() {
	u.router.Post("/token", handlers.WrapError(u.token, u.logger))
}

func (u *User) token(w http.ResponseWriter, r *http.Request) error {
	var body GetAccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return utils.Err400(w, err, "Error decoding body", nil)
	}

	req, err := body.ToUseCase()
	if err != nil {
		return utils.Err400(w, err, err.Error(), nil)
	}
	resUC, errUC := u.userUseCase.GetAccessToken(req)
	if errUC != nil {
		return errUC.SendError(w)
	}

	res := GetAccessTokenResponse{
		AccessToken: resUC.Token.Token,
		ExpireAt:    resUC.Token.ExpiredAt,
	}

	return utils.JSON(w, res)
}
