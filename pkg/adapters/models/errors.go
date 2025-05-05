package models

import "errors"

// Models errors list
var (
	ErrIDFromString       = errors.New("error when a new ID from a string")
	ErrPasswordFromString = errors.New("error when a new password from a string")
	ErrEmailFromString    = errors.New("error when a new email from a string")
	ErrParseDateTime      = errors.New("error when parsing date time")
)
