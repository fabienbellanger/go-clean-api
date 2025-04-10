package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
