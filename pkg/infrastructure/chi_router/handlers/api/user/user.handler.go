package user

import (
	"encoding/json"
	"errors"
	domainerr "go-clean-api/pkg/domain/errors"
	"go-clean-api/pkg/domain/usecases"
	vo "go-clean-api/pkg/domain/value_objects"
	"go-clean-api/pkg/infrastructure/chi_router/handlers"
	"go-clean-api/pkg/infrastructure/chi_router/httputil"
	"go-clean-api/pkg/infrastructure/logger"
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
	h.router.Get("/", handlers.WrapError(h.GetAll, h.logger))
	h.router.Get("/deleted", handlers.WrapError(h.GetAllDeleted, h.logger))
	h.router.Get("/{id}", handlers.WrapError(h.getByID, h.logger))
	h.router.Delete("/{id}", handlers.WrapError(h.delete, h.logger))
	h.router.Patch("/{id}/restore", handlers.WrapError(h.restore, h.logger))
}

func (u *Handler) token(w http.ResponseWriter, r *http.Request) error {
	var body GetAccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httputil.Err400(w, err, "Error when decoding the body", nil)
	}

	req, err := body.ToUseCase()
	if err != nil {
		return httputil.Err400(w, err, "Invalid parameters", err)
	}

	resUC, errUC := u.userUseCase.GetAccessToken(req)
	if errUC != nil {
		if errors.Is(errUC, domainerr.ErrNotFound) || errors.Is(errUC, usecases.ErrInvalidPassword) {
			return httputil.Err401(w, errUC, "Unauthorized", nil)
		} else if errors.Is(errUC, usecases.ErrAccessTokenCreation) {
			return httputil.Err500(w, errUC, "Internal server error", "Error during token generation")
		} else if errors.Is(errUC, domainerr.ErrDatabase) {
			return httputil.Err500(w, errUC, "Internal server error", "Error during authentication")
		} else {
			return httputil.Err500(w, errUC, "Internal server error", "Unknown error")
		}
	}

	res := GetAccessTokenResponse{
		AccessToken:          resUC.Token.Token,
		AccessTokenExpiredAt: resUC.Token.ExpiredAt.RFC3339(),
	}

	return httputil.JSON(w, res)
}

func (u *Handler) register(w http.ResponseWriter, r *http.Request) error {
	var body CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httputil.Err400(w, err, "Error when decoding the body", nil)
	}

	req, err := body.ToUseCase()
	if err != nil {
		return httputil.Err400(w, err, "Invalid parameters", err)
	}

	resUC, errUC := u.userUseCase.Create(req)
	if errUC != nil {
		return httputil.Err500(w, errUC, "Internal server error", "Error during user creation")
	}

	res := CreateResponse{
		ID:        resUC.ID.String(),
		Email:     resUC.Email.Value(),
		Lastname:  resUC.Lastname,
		Firstname: resUC.Firstname,
		CreatedAt: resUC.CreatedAt.RFC3339(),
		UpdatedAt: resUC.UpdatedAt.RFC3339(),
	}

	return httputil.Created(w, res)
}

func (u *Handler) getByID(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return httputil.Err400(w, nil, "ID is required", nil)
	}

	req, err := GetByIDRequest{ID: id}.ToUseCase()
	if err != nil {
		return httputil.Err400(w, err, "Invalid parameters", err)
	}

	resUC, errUC := u.userUseCase.GetByID(req)
	if errUC != nil {
		if errors.Is(errUC, domainerr.ErrNotFound) {
			return httputil.Err404(w, errUC, "No user found", nil)
		} else if errors.Is(errUC, domainerr.ErrDatabase) {
			return httputil.Err500(w, errUC, "Internal server error", "Error when getting user")
		} else {
			return httputil.Err500(w, errUC, "Internal server error", "Unknown error")
		}
	}

	res := GetByIDResponse{}.FromEntity(resUC)

	return httputil.JSON(w, res)
}

func (u *Handler) GetAll(w http.ResponseWriter, r *http.Request) error {
	p := r.URL.Query().Get("page")
	s := r.URL.Query().Get("size")
	pagination := vo.PaginationFromQuery(p, s, "")

	users, errUC := u.userUseCase.GetAll(usecases.GetAllUsersRequest{
		Pagination: pagination,
		Deleted:    false,
	})
	if errUC != nil {
		return httputil.Err500(w, errUC, "Internal server error", "Error when getting users")
	}

	res := GetAllResponse{}.FromEntity(users, pagination)

	return httputil.JSON(w, res)
}

func (u *Handler) GetAllDeleted(w http.ResponseWriter, r *http.Request) error {
	p := r.URL.Query().Get("page")
	s := r.URL.Query().Get("size")
	pagination := vo.PaginationFromQuery(p, s, "")

	users, errUC := u.userUseCase.GetAll(usecases.GetAllUsersRequest{
		Pagination: pagination,
		Deleted:    true,
	})
	if errUC != nil {
		return httputil.Err500(w, errUC, "Internal server error", "Error when getting deleted users")
	}

	res := GetAllResponse{}.FromEntity(users, pagination)

	return httputil.JSON(w, res)
}

func (u *Handler) delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return httputil.Err400(w, nil, "ID is required", nil)
	}

	req, err := DeleteRestoreRequest{ID: id}.ToUseCase()
	if err != nil {
		return httputil.Err400(w, err, "Invalid parameters", err)
	}

	_, errUC := u.userUseCase.Delete(req)
	if errUC != nil {
		if errors.Is(errUC, domainerr.ErrNotFound) {
			return httputil.Err404(w, errUC, "No user found", nil)
		} else if errors.Is(errUC, domainerr.ErrDatabase) {
			return httputil.Err500(w, errUC, "Internal server error", "Error when deleting user")
		} else {
			return httputil.Err500(w, errUC, "Internal server error", "Unknown error")
		}
	}

	return httputil.NoContent(w)
}

func (u *Handler) restore(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return httputil.Err400(w, nil, "ID is required", nil)
	}

	req, err := DeleteRestoreRequest{ID: id}.ToUseCase()
	if err != nil {
		return httputil.Err400(w, err, "Invalid parameters", err)
	}

	_, errUC := u.userUseCase.Restore(req)
	if errUC != nil {
		if errors.Is(errUC, domainerr.ErrNotFound) {
			return httputil.Err404(w, errUC, "No user found", nil)
		} else if errors.Is(errUC, domainerr.ErrDatabase) {
			return httputil.Err500(w, errUC, "Internal server error", "Error when restoring user")
		} else {
			return httputil.Err500(w, errUC, "Internal server error", "Unknown error")
		}
	}

	return httputil.NoContent(w)
}
