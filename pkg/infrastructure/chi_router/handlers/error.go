package handlers

import (
	"fmt"
	"go-clean-api/pkg/infrastructure/logger"
	"net/http"
)

// RequestIDKey is the key used to store the request ID in the context
type RequestIDKey string

// WrapError wraps the handlers error and logs it
func WrapError(f func(w http.ResponseWriter, r *http.Request) error, l logger.CustomLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := fmt.Sprintf("%s", r.Context().Value(RequestIDKey("request_id")))

		err := f(w, r)
		if err != nil {
			fields := logger.Fields{
				logger.NewField("request_id", "string", requestId),
			}
			l.Error(err.Error(), fields)
		}
	}
}
