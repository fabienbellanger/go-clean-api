package values_objects

import "strconv"

const (
	// Pagination min size
	PaginationMinSize = 50

	// Pagination min size
	PaginationMaxSize = 500

	// Pagination min size
	PaginationDefaultSize = 100
)

type Pagination struct {
	page    int
	size    int
	maxSize int
}

// NewPagination creates a new Pagination
func NewPagination(page, size, maxSize int) Pagination {
	p := max(page, 1)

	m := maxSize
	if maxSize == 0 || maxSize > PaginationMaxSize {
		m = PaginationMaxSize
	}

	s := size
	if size > m {
		s = m
	} else if size < PaginationMinSize {
		s = PaginationMinSize
	}

	return Pagination{
		page:    p,
		size:    s,
		maxSize: m,
	}
}

// TODO: Add test
func PaginationFromQuery(page, size, maxSize string) Pagination {
	p, err := strconv.Atoi(page)
	if err != nil || p < 1 {
		p = 1
	}

	m, err := strconv.Atoi(maxSize)
	if err != nil || m == 0 || m > PaginationMaxSize {
		m = PaginationMaxSize
	}

	s, err := strconv.Atoi(size)
	if err != nil {
		s = PaginationDefaultSize
	} else if s > m {
		s = m
	} else if s < PaginationMinSize {
		s = PaginationMinSize
	}

	return Pagination{
		page:    p,
		size:    s,
		maxSize: m,
	}
}

func (p Pagination) Page() int {
	return p.page
}

func (p Pagination) Size() int {
	return p.size
}
