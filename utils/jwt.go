package utils

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// LoadECDSAKeyFromFile loads an ECDSA private or public key from a file
func LoadECDSAKeyFromFile(filename string, isPrivate bool) (any, error) {
	// Read file
	pemBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, NewAppErr(err, fmt.Sprintf("error when reading file: %s", filename), nil, nil)
	}

	// Decode PEM
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, NewAppErr(
			errors.New("error when decoding .pem file"),
			fmt.Sprintf("error when decoding %s file", filename),
			nil,
			nil,
		)
	}

	// Parse key
	var key any
	if isPrivate {
		key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	} else {
		key, err = x509.ParsePKIXPublicKey(block.Bytes)
	}
	if err != nil {
		return nil, NewAppErr(err, "error when parsing JWT key", nil, nil)
	}

	return key, nil
}

// GetTokenAndKeyFromAlgo returns a token and a key from an algorithm and a secret
func GetTokenAndKeyFromAlgo(algo, secret, keyPath string) (*jwt.Token, any, error) {
	// Create token
	var token *jwt.Token
	var key any
	var err error

	if algo == "HS512" {
		if len(secret) < 8 {
			return nil, nil, NewAppErr(
				errors.New("secret must have at least 8 characters"),
				"secret must have at least 8 characters",
				nil,
				nil,
			)
		}

		token = jwt.New(jwt.SigningMethodHS512)

		key = []byte(secret)
	} else if algo == "ES384" {
		token = jwt.New(jwt.SigningMethodES384)

		key, err = LoadECDSAKeyFromFile(keyPath, true)
		if err != nil {
			return nil, nil, NewAppErr(err, "error when loading an ECDSA private or public key from a file", nil, nil)
		}
	} else {
		return nil, nil, NewAppErr(
			errors.New("unsupported JWT algo: must be HS512 or ES384"),
			"unsupported JWT algo: must be HS512 or ES384",
			nil,
			nil,
		)
	}

	return token, key, nil
}

// GetKeyFromAlgo returns a key from an algorithm and a secret
func GetKeyFromAlgo(algo, secret, keyPath string) (any, error) {
	var key any
	var err error

	if algo == jwt.SigningMethodHS512.Name {
		key = []byte(secret)
	} else if algo == jwt.SigningMethodES384.Name {
		key, err = LoadECDSAKeyFromFile(keyPath, false)
		if err != nil {
			return nil, NewAppErr(err, "error when loading an ECDSA private or public key from a file", nil, nil)
		}
	} else {
		return nil, NewAppErr(
			errors.New("unsupported JWT algo: must be HS512 or ES384"),
			"unsupported JWT algo: must be HS512 or ES384",
			nil,
			nil,
		)
	}

	return key, nil
}
