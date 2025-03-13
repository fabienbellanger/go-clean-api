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
	h.router.Post("/", handlers.WrapError(h.register, h.logger))
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
		AccessToken:          resUC.Token.Token,
		AccessTokenExpiredAt: resUC.Token.ExpiredAt,
	}

	return utils.JSON(w, res)
}

func (u *Handler) register(w http.ResponseWriter, r *http.Request) error {
	var body CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return utils.Err400(w, err, "Error when decoding the body", nil)
	}

	req, err := body.ToUseCase()
	if err != nil {
		return utils.Err400(w, err, "Invalid parameters", err)
	}

	resUC, errUC := u.userUseCase.Create(req)
	if errUC != nil {
		return errUC.SendError(w)
	}

	res := CreateResponse{
		ID:        resUC.ID.String(),
		Email:     resUC.Email.Value(),
		Lastname:  resUC.Lastname,
		Firstname: resUC.Firstname,
		CreatedAt: utils.FormatToRFC3339(resUC.CreatedAt),
		UpdatedAt: utils.FormatToRFC3339(resUC.UpdatedAt),
	}

	return utils.JSON(w, res)
}
