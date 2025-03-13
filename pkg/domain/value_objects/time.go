package values_objects

import (
	"errors"
	"time"
)

const sqlFormat = "2006-01-02 15:04:05"

// Time represents an time value object
type Time struct {
	value time.Time
}

// NewTime creates a new time
//
// The default timezone is UTC.
func NewTime(value time.Time, loc *time.Location) Time {
	if loc == nil {
		loc = time.UTC
	}

	// Set location
	value = value.In(loc)

	return Time{value: value}
}

// Value returns the time value
func (e Time) Value() time.Time {
	return e.value
}

// SQL returns the time in SQL format
func (d Time) SQL() string {
	return d.Value().Format(sqlFormat)
}

// RFC3339 returns the time in SQL format
func (d Time) RFC3339() string {
	return d.Value().Format(time.RFC3339)
}

// ParseSQL parses a time from SQL format
func ParseSQL(s string, loc *time.Location) (Time, error) {
	if loc == nil {
		return Time{}, errors.New("location is required")
	}

	t, err := time.ParseInLocation(sqlFormat, s, loc)
	if err != nil {
		return Time{}, err
	}

	return NewTime(t, loc), nil
}

// ParseRFC3339 parses a time from RFC3339 format
func ParseRFC3339(s string, loc *time.Location) (Time, error) {
	t, err := time.ParseInLocation(time.RFC3339, s, loc)
	if err != nil {
		return Time{}, err
	}

	return NewTime(t, loc), nil
}
