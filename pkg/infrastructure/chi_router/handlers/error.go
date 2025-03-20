package handlers

import (
	"fmt"
	"go-clean-api/pkg/infrastructure/logger"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestIDKey is the key used to store the request ID in the context
type RequestIDKey string

// WrapError wraps the handlers error and logs it
func WrapError(f func(w http.ResponseWriter, r *http.Request) error, l logger.CustomLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := fmt.Sprintf("%s", r.Context().Value(RequestIDKey("request_id")))
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		err := f(ww, r)

		if err != nil {
			fields := logger.Fields{
				logger.NewField("request_id", "string", requestId),
			}

			if ww.Status() == http.StatusInternalServerError {
				l.Error(err.Error(), fields)
			} else {
				l.Warn(err.Error(), fields)
			}
		}
	}
}
