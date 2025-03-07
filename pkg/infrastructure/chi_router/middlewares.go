package chi_router

import (
	"context"
	"fmt"
	"go-clean-api/pkg/infrastructure/chi_router/handlers"
	"go-clean-api/pkg/infrastructure/logger"
	"go-clean-api/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func (s *ChiServer) initJWTToken() error {
	algo := s.Config.JWT.Algorithm
	key, err := utils.GetKeyFromAlgo(algo, s.Config.JWT.PrivateKeyPath, s.Config.JWT.PublicKeyPath)
	if err != nil {
		return err
	}

	tokenAuth = jwtauth.New(algo, key, nil)

	return nil
}

func (s *ChiServer) initMiddlewares(r *chi.Mux) {
	r.Use(s.requestID) // Must be before the access logger
	if s.Config.Log.EnableAccessLog {
		r.Use(s.initAccessLogger())
	}

	// Request size limiter
	if s.Config.Server.MaxRequestSize > 0 {
		r.Use(middleware.RequestSize(s.Config.Server.MaxRequestSize << 10))
	}

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(s.Config.Server.Timeout) * time.Second))
	r.Use(middleware.RealIP)

	// Profiler
	if s.Config.Pprof.Enable {
		r.Group(func(r chi.Router) {
			// TODO: Change?
			// r.Use(s.initPprofBasicAuth())

			r.Mount("/debug", middleware.Profiler())
		})
	}
}

func (s *ChiServer) initAccessLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now().UTC()
			ctxID := r.Context().Value(handlers.RequestIDKey("request_id"))
			var requestId string
			if ctxID != nil {
				requestId = fmt.Sprintf("%s", r.Context().Value(handlers.RequestIDKey("request_id")))
			}
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			stop := time.Since(start)
			url := r.Host + r.RequestURI // TODO: Do better, missing https:// or http://
			fields := logger.Fields{
				logger.NewField("code", "int", ww.Status()),
				logger.NewField("method", "string", r.Method),
				logger.NewField("path", "string", r.URL.Path),
				logger.NewField("url", "string", url),
				logger.NewField("ip", "string", r.RemoteAddr), // TODO: Remove port
				logger.NewField("userAgent", "string", r.UserAgent()),
				logger.NewField("latency", "string", stop.String()),
				logger.NewField("request_id", "string", requestId),
			}

			s.Logger.Info("", fields)
		}
		return http.HandlerFunc(fn)
	}
}

func (s *ChiServer) initCORS() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   s.Config.CORS.AllowedOrigins,
		AllowedMethods:   s.Config.CORS.AllowedMethods,
		AllowedHeaders:   s.Config.CORS.AllowedHeaders,
		ExposedHeaders:   s.Config.CORS.ExposedHeaders,
		AllowCredentials: s.Config.CORS.AllowCredentials,
		MaxAge:           s.Config.CORS.MaxAge,
	})
}

func (s *ChiServer) initBasicAuth() func(next http.Handler) http.Handler {
	creds := make(map[string]string, 1)
	creds[s.Config.Server.BasicAuthUsername] = s.Config.Server.BasicAuthPassword

	return middleware.BasicAuth("Restricted", creds)
}

func (s *ChiServer) initPprofBasicAuth() func(next http.Handler) http.Handler {
	creds := make(map[string]string, 1)
	creds[s.Config.Pprof.BasicAuthUsername] = s.Config.Pprof.BasicAuthPassword

	return middleware.BasicAuth("Restricted", creds)
}

func (s *ChiServer) initJWT(r chi.Router) {
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(s.jwtAuthenticator(tokenAuth))
}

func (s *ChiServer) jwtAuthenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				utils.Err401(w, err, "Unauthorized", nil) // TODO: Error not managed
				return
			}

			if token == nil || jwt.Validate(token, ja.ValidateOptions()...) != nil {
				utils.Err401(w, nil, "Unauthorized", nil) // TODO: Error not managed
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

func (s *ChiServer) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), handlers.RequestIDKey("request_id"), id)

		w.Header().Add("X-Request-Id", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
