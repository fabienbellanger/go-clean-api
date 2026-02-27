package db

import (
	"fmt"
	vo "go-clean-api/pkg/domain/value_objects"
	"strings"
	"time"
)

type DB interface {
	DSN() (string, error)
	Database(string)
}

const (
	// DefaultSlowThreshold represents the default slow threshold value
	DefaultSlowThreshold time.Duration = 200 * time.Millisecond
)

// PaginateValues transforms page and limit into offset and limit.
func PaginateValues(p, l int) (offset int, limit int) {
	if p < 1 {
		p = 1
	}

	limit = l
	if limit > vo.PaginationMaxSize || limit < 1 {
		limit = vo.PaginationMaxSize
	}

	offset = (p - 1) * limit

	return
}

// orderValues transforms list of fields to sort into a map.
func orderValues(list string, prefixes ...string) []string {
	r := make([]string, 0)

	if len(list) <= 0 {
		return r
	}

	prefix := ""
	if len(prefixes) == 1 {
		prefix = prefixes[0] + "."
	}

	sorts := strings.SplitSeq(list, ",")
	for s := range sorts {
		if len(s) > 0 {
			key := fmt.Sprintf("%s%s", prefix, s[1:])
			if strings.HasPrefix(s, "+") && len(s[1:]) > 1 {
				r = append(r, fmt.Sprintf("%s ASC", key))
			} else if strings.HasPrefix(s, "-") && len(s[1:]) > 1 {
				r = append(r, fmt.Sprintf("%s DESC", key))
			}
		}
	}

	return r
}

// OrderValues returns the ORDER BY clause for a list of fields to sort.
func OrderValues(list string, prefixes ...string) (s string) {
	values := orderValues(list, prefixes...)
	s = strings.Join(values, ", ")

	if len(s) > 0 {
		s = " ORDER BY " + s
	}

	return
}
