package values_objects

const (
	// Pagination min size
	PaginationMinSize = 50

	// Pagination min size
	PaginationMaxSize = 500

	// Pagination min size
	PaginationDefaultSize = 100
)

type Pagination struct {
	page    uint
	size    uint
	maxSize uint
}

// NewPagination creates a new Pagination
func NewPagination(page, size, maxSize uint) Pagination {
	p := page
	if page == 0 {
		p = 1
	}

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

func (p Pagination) Page() uint {
	return p.page
}

func (p Pagination) Size() uint {
	return p.size
}
