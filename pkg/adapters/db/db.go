package db

import (
	"errors"
	"fmt"
	"go-clean-api/pkg"
	values_objects "go-clean-api/pkg/domain/value_objects"
	"strings"
	"time"
)

const (
	// DefaultSlowThreshold represents the default slow threshold value
	DefaultSlowThreshold time.Duration = 200 * time.Millisecond
)

// Config represents the MySQL database configuration
type Config struct {
	pkg.ConfigDatabase
}

// dsn returns the DSN if the configuration is OK or an error in other case
func (c *Config) dsn() (dsn string, err error) {
	if c.Host == "" || c.Port == 0 || c.Username == "" || c.Password == "" {
		return dsn, errors.New("error in database configuration")
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database)
	if c.Charset != "" {
		dsn += fmt.Sprintf("&charset=%s", c.Charset)
	}
	if c.Collation != "" {
		dsn += fmt.Sprintf("&collation=%s", c.Collation)
	}
	if c.Location != "" {
		dsn += fmt.Sprintf("&loc=%s", c.Location)
	}
	return
}

// PaginateValues transforms page and limit into offset and limit.
func PaginateValues(p, l int) (offset int, limit int) {
	if p < 1 {
		p = 1
	}

	limit = l
	if limit > values_objects.PaginationMaxSize || limit < 1 {
		limit = values_objects.PaginationMaxSize
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

	sorts := strings.Split(list, ",")
	for _, s := range sorts {
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
