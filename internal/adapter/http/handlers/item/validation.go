package item

import (
	"net/url"
	"strconv"

	"github.com/1337jazz/generoni/internal/core/domain"
)

func validateGetItem(rawID string) (int64, error) {
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id < 1 {
		return 0, domain.ErrInvalidItemID
	}
	return id, nil
}

func validateListItems(query url.Values) (domain.Pagination, error) {
	limit, err := parseOptionalInt(query.Get("limit"))
	if err != nil {
		return domain.Pagination{}, domain.ErrInvalidLimit
	}
	offset, err := parseOptionalInt(query.Get("offset"))
	if err != nil {
		return domain.Pagination{}, domain.ErrInvalidOffset
	}
	return domain.NewPagination(limit, offset)
}

func parseOptionalInt(raw string) (int64, error) {
	if raw == "" {
		return 0, nil
	}
	return strconv.ParseInt(raw, 10, 64)
}
