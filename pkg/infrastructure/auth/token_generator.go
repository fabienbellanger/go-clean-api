package auth

import (
	"time"

	"go-clean-api/pkg"
	"go-clean-api/pkg/apperr"
	"go-clean-api/pkg/domain/entities"
	vo "go-clean-api/pkg/domain/value_objects"

	"github.com/golang-jwt/jwt/v5"
)

// JWTTokenGenerator implements services.TokenGenerator using JWT.
type JWTTokenGenerator struct {
	cfg pkg.ConfigJWT
}

// NewJWTTokenGenerator creates a new JWTTokenGenerator.
func NewJWTTokenGenerator(cfg pkg.ConfigJWT) *JWTTokenGenerator {
	return &JWTTokenGenerator{cfg: cfg}
}

// Generate creates a new JWT access token for the given user ID.
func (g *JWTTokenGenerator) Generate(id entities.UserID) (entities.AccessToken, error) {
	// Create token and key
	token, key, err := GetTokenAndKeyFromAlgo(g.cfg.Algorithm, g.cfg.SecretKey, g.cfg.PrivateKeyPath)
	if err != nil {
		return entities.AccessToken{}, err
	}

	// Expiration time
	now := time.Now()
	expiredAt := now.Add(g.cfg.Lifetime)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = id.String()
	claims["exp"] = expiredAt.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	// Generate encoded token and send it as response
	t, err := token.SignedString(key)
	if err != nil {
		return entities.AccessToken{}, apperr.NewAppErr(
			err,
			"error when signing JWT token",
			nil,
			nil,
		)
	}
	return entities.AccessToken{Token: t, ExpiredAt: vo.NewTime(expiredAt, nil)}, nil
}
