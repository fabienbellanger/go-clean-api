package user

import (
	"encoding/json"
	"go-clean-api/pkg/domain/usecases"
	"go-clean-api/pkg/infrastructure/chi_router/handlers"
	"go-clean-api/pkg/infrastructure/logger"
	"go-clean-api/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler handles user requests
type Handler struct {
	router      chi.Router
	userUseCase usecases.User
	logger      logger.CustomLogger
}

// NewHandler returns a new Handler
func NewHandler(r chi.Router, l logger.CustomLogger, userUseCase usecases.User) Handler {
	return Handler{
		router:      r,
		userUseCase: userUseCase,
		logger:      l,
	}
}

// PublicRoutes adds users public routes
func (h *Handler) PublicRoutes() {
	h.router.Post("/token", handlers.WrapError(h.token, h.logger))
}

// PrivateRoutes adds users private routes
func (h *Handler) PrivateRoutes() {
}

func (u *Handler) token(w http.ResponseWriter, r *http.Request) error {
	var body GetAccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return utils.Err400(w, err, "Error when decoding the body", nil)
	}

	req, err := body.ToUseCase()
	if err != nil {
		return utils.Err400(w, err, "Invalid parameters", err)
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
