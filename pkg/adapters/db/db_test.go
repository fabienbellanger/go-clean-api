package db

import (
	"errors"
	"go-clean-api/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDsn(t *testing.T) {
	type result struct {
		dsn string
		err error
	}

	tests := []struct {
		name   string
		args   Config
		wanted result
	}{
		{
			name: "Simple valid DSN",
			args: Config{
				ConfigDatabase: pkg.ConfigDatabase{
					Username: "root",
					Password: "root",
					Database: "test",
					Host:     "localhost",
					Port:     3306,
				},
			},
			wanted: result{
				dsn: "root:root@tcp(localhost:3306)/test?parseTime=True",
				err: nil,
			},
		},
		{
			name: "Complet valid DSN",
			args: Config{
				ConfigDatabase: pkg.ConfigDatabase{
					Username:  "root",
					Password:  "root",
					Database:  "test",
					Host:      "localhost",
					Port:      3306,
					Charset:   "utf8mb4",
					Collation: "utf8mb4_general_ci",
					Location:  "Local",
				},
			},
			wanted: result{
				dsn: "root:root@tcp(localhost:3306)/test?parseTime=True&charset=utf8mb4&collation=utf8mb4_general_ci&loc=Local",
				err: nil,
			},
		},
		{
			name: "Invalid DSN (no username)",
			args: Config{
				ConfigDatabase: pkg.ConfigDatabase{
					Password: "root",
					Database: "test",
					Port:     3306,
					Host:     "localhost",
				},
			},
			wanted: result{
				dsn: "",
				err: errors.New("error in database configuration"),
			},
		},
		{
			name: "Invalid DSN (no password)",
			args: Config{
				ConfigDatabase: pkg.ConfigDatabase{
					Username: "root",
					Database: "test",
					Port:     3306,
					Host:     "localhost",
				},
			},
			wanted: result{
				dsn: "",
				err: errors.New("error in database configuration"),
			},
		},
		{
			name: "Invalid DSN (no port)",
			args: Config{
				ConfigDatabase: pkg.ConfigDatabase{
					Username: "root",
					Password: "root",
					Database: "test",
					Host:     "localhost",
				},
			},
			wanted: result{
				dsn: "",
				err: errors.New("error in database configuration"),
			},
		},
		{
			name: "Invalid DSN (no host)",
			args: Config{
				ConfigDatabase: pkg.ConfigDatabase{
					Username: "root",
					Password: "root",
					Database: "test",
					Port:     3306,
				},
			},
			wanted: result{
				dsn: "",
				err: errors.New("error in database configuration"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn, err := tt.args.dsn()
			got := result{dsn, err}

			if got.err != nil {
				assert.Equal(t, got.dsn, tt.wanted.dsn)
			}
			assert.Equal(t, got.err, tt.wanted.err)
		})
	}
}

func TestPaginateValues(t *testing.T) {
	type result struct {
		offset int
		limit  int
	}

	tests := []struct {
		name   string
		args   []int
		wanted result
	}{
		{
			name: "Simple valid pagination",
			args: []int{1, 10},
			wanted: result{
				offset: 0,
				limit:  10,
			},
		},
		{
			name: "Invalid page",
			args: []int{0, 10},
			wanted: result{
				offset: 0,
				limit:  10,
			},
		},
		{
			name: "Invalid limit",
			args: []int{1, 0},
			wanted: result{
				offset: 0,
				limit:  500,
			},
		},
		{
			name: "Invalid page and limit",
			args: []int{0, 0},
			wanted: result{
				offset: 0,
				limit:  500,
			},
		},
		{
			name: "Limit too high",
			args: []int{1, 200},
			wanted: result{
				offset: 0,
				limit:  200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset, limit := PaginateValues(tt.args[0], tt.args[1])
			got := result{offset, limit}

			assert.Equal(t, got, tt.wanted)
		})
	}
}

func TestOrderValues(t *testing.T) {
	type result struct {
		sort string
	}

	tests := []struct {
		name   string
		args   []string
		wanted result
	}{
		{
			name: "Simple sort",
			args: []string{"+id"},
			wanted: result{
				sort: " ORDER BY id ASC",
			},
		},
		{
			name: "Many filed",
			args: []string{"+id,-name,+created_at"},
			wanted: result{
				sort: " ORDER BY id ASC, name DESC, created_at ASC",
			},
		},
		{
			name: "Empty",
			args: []string{""},
			wanted: result{
				sort: "",
			},
		},
		{
			name: "One invalid field",
			args: []string{"+id,name,+created_at"},
			wanted: result{
				sort: " ORDER BY id ASC, created_at ASC",
			},
		},
		{
			name: "With prefix",
			args: []string{"+id,name,+created_at", "users"},
			wanted: result{
				sort: " ORDER BY users.id ASC, users.created_at ASC",
			},
		},
		{
			name: "With prefix and all fields invalid",
			args: []string{"id,name;created_a", "users"},
			wanted: result{
				sort: "",
			},
		},
		{
			name: "With empty last field",
			args: []string{"-id,+name,"},
			wanted: result{
				sort: " ORDER BY id DESC, name ASC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sort string
			if len(tt.args) == 1 {
				sort = OrderValues(tt.args[0])
			} else if len(tt.args) == 2 {
				sort = OrderValues(tt.args[0], tt.args[1])
			}
			got := result{sort}

			assert.Equal(t, got, tt.wanted)
		})
	}
}
