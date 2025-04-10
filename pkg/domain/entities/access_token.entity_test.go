package entities

import (
	"errors"
	"go-clean-api/pkg"
	vo "go-clean-api/pkg/domain/value_objects"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	type args struct {
		userID   UserID
		lifetime time.Duration
		algo     string
		secret   string
	}

	type result struct {
		jwt AccessToken
		err error
	}

	lifetime := time.Duration(2) * time.Hour

	tests := []struct {
		name   string
		args   args
		wanted result
	}{
		{
			name: "Invalid algo",
			args: args{
				userID:   vo.NewID(),
				lifetime: lifetime,
				algo:     "",
				secret:   "my-secret",
			},
			wanted: result{
				jwt: AccessToken{},
				err: errors.New("unsupported JWT algo: must be HS512 or ES384"),
			},
		},
		{
			name: "Invalid algo",
			args: args{
				userID:   vo.NewID(),
				lifetime: lifetime,
				algo:     "HS512",
				secret:   "secret",
			},
			wanted: result{
				jwt: AccessToken{},
				err: errors.New("secret must have at least 8 characters"),
			},
		},
		{
			name: "Valid",
			args: args{
				userID:   vo.NewID(),
				lifetime: lifetime,
				algo:     "HS512",
				secret:   "my-secret",
			},
			wanted: result{
				jwt: AccessToken{},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := pkg.ConfigJWT{
				Algorithm: tt.args.algo,
				Lifetime:  tt.args.lifetime,
				SecretKey: tt.args.secret,
			}
			jwt, err := NewAccessToken(
				tt.args.userID,
				cfg,
			)
			got := result{jwt, err}

			if got.err != nil {
				assert.Equal(t, got.jwt.Token, tt.wanted.jwt.Token)
			} else {
				assert.Greater(t, len(got.jwt.Token), 0)
				assert.Greater(t, got.jwt.ExpiredAt.Value(), time.Now().Add(lifetime-time.Minute))
				assert.Less(t, got.jwt.ExpiredAt.Value(), time.Now().Add(lifetime+time.Minute))
			}
			assert.Equal(t, got.err, tt.wanted.err)
		})
	}
}
