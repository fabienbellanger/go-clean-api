package utils

import "github.com/fabienbellanger/xerr"

type AppErr = xerr.Err

// NewAppErr creates a new AppErr instance with the provided error, message, details, and previous error
func NewAppErr(err error, message string, details any, prev *AppErr) AppErr {
	return xerr.New(err, message, details, 0, prev)
}

// NewAppErrWithCode is a custom error type that extends the xerr package with a specific error code
func NewAppErrWithCode(err error, message string, details any, code int, prev *AppErr) AppErr {
	return xerr.New(err, message, details, code, prev)
}

// EmptyAppErr returns an empty AppErr
func EmptyAppErr() AppErr {
	return xerr.Empty()
}
