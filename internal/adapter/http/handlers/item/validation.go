package item

import (
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
