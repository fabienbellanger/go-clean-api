package values_objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name    string
		page    uint
		size    uint
		maxSize uint
		want    Pagination
	}{
		{
			name:    "Test NewPagination with valid values",
			page:    1,
			size:    100,
			maxSize: 0,
			want: Pagination{
				page:    1,
				size:    100,
				maxSize: PaginationMaxSize,
			},
		},
		{
			name:    "Test NewPagination with invalid page",
			page:    0,
			size:    100,
			maxSize: 0,
			want: Pagination{
				page:    1,
				size:    100,
				maxSize: PaginationMaxSize,
			},
		},
		{
			name:    "Test NewPagination with size too small",
			page:    2,
			size:    10,
			maxSize: 0,
			want: Pagination{
				page:    2,
				size:    PaginationMinSize,
				maxSize: PaginationMaxSize,
			},
		},
		{
			name:    "Test NewPagination with size too big",
			page:    2,
			size:    1_000,
			maxSize: 0,
			want: Pagination{
				page:    2,
				size:    PaginationMaxSize,
				maxSize: PaginationMaxSize,
			},
		},
		{
			name:    "Test NewPagination with max size",
			page:    2,
			size:    100,
			maxSize: 400,
			want: Pagination{
				page:    2,
				size:    100,
				maxSize: 400,
			},
		},
		{
			name:    "Test NewPagination with size too big and max size",
			page:    2,
			size:    1_000,
			maxSize: 600,
			want: Pagination{
				page:    2,
				size:    PaginationMaxSize,
				maxSize: PaginationMaxSize,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewPagination(tt.page, tt.size, tt.maxSize), tt.want)
		})
	}
}
