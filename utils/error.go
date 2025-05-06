package utils

import "github.com/fabienbellanger/xerr"

type AppErr = xerr.Err

// AppErr is a custom error type that extends the xerr package
func NewAppErr(err error, message string, details any, prev *AppErr) AppErr {
	return xerr.NewErr(err, message, details, prev)
}

// NewAppErrFromCode creates a new AppErr with a specific code
func EmptyAppErr() AppErr {
	return xerr.EmptyErr()
}
