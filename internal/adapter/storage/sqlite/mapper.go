package sqlite

import (
	"github.com/1337jazz/generoni/internal/adapter/storage/sqlite/queries"
	"github.com/1337jazz/generoni/internal/core/domain"
)

func itemRowToDomain(row queries.Item) *domain.Item {
	return &domain.Item{
		ID:        row.ID,
		Title:     row.Title,
		CreatedAt: row.CreatedAt,
	}
}
