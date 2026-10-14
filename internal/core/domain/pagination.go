package domain

const (
	DefaultPageLimit = 25
	MaxPageLimit     = 100
)

type Pagination struct {
	Limit  int64
	Offset int64
}

func NewPagination(limit, offset int64) (Pagination, error) {
	if limit < 0 || limit > MaxPageLimit {
		return Pagination{}, ErrInvalidLimit
	}
	if offset < 0 {
		return Pagination{}, ErrInvalidOffset
	}
	if limit == 0 {
		limit = DefaultPageLimit
	}
	return Pagination{Limit: limit, Offset: offset}, nil
}
