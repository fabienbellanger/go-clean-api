package apperr

import "github.com/fabienbellanger/xerr"

type AppErr = xerr.Err

// NewAppErr creates a new AppErr instance with the provided error, message, details, and previous error
func NewAppErr(err error, message string, details any, prev *AppErr) *AppErr {
	e := xerr.New(err, message, details, 0, prev, 2)
	return &e
}

// NewAppErrWithCode is a custom error type that extends the xerr package with a specific error code
func NewAppErrWithCode(err error, message string, details any, code int, prev *AppErr) *AppErr {
	e := xerr.New(err, message, details, code, prev, 2)
	return &e
}

// EmptyAppErr returns an empty AppErr
func EmptyAppErr() *AppErr {
	e := xerr.Empty()
	return &e
}
