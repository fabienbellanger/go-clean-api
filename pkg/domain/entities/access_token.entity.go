package entities

import (
	"go-clean-api/pkg"
	vo "go-clean-api/pkg/domain/value_objects"
	"go-clean-api/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessToken is a struct that represents a JWT access token
type AccessToken struct {
	Token     string
	ExpiredAt vo.Time
}

func NewAccessToken(id UserID, cfg pkg.ConfigJWT) (AccessToken, error) {
	// Create token and key
	token, key, err := utils.GetTokenAndKeyFromAlgo(cfg.Algorithm, cfg.SecretKey, cfg.PrivateKeyPath)
	if err != nil {
		return AccessToken{}, err
	}

	// Expiration time
	now := time.Now()
	expiredAt := now.Add(cfg.Lifetime)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = id.String()
	claims["exp"] = expiredAt.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	// Generate encoded token and send it as response
	t, err := token.SignedString(key)
	if err != nil {
		return AccessToken{}, err
	}
	return AccessToken{Token: t, ExpiredAt: vo.NewTime(expiredAt, nil)}, nil
}
