package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/1337jazz/generoni/internal/adapter/storage/sqlite/queries"
	"github.com/1337jazz/generoni/internal/core/domain"
	"github.com/1337jazz/generoni/internal/core/port"
)

var _ port.ItemRepository = (*ItemRepository)(nil)

type ItemRepository struct {
	db      *sql.DB
	queries *queries.Queries
}

func NewItemRepository(db *sql.DB, queries *queries.Queries) *ItemRepository {
	return &ItemRepository{db: db, queries: queries}
}

func (r *ItemRepository) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	row, err := r.queries.CreateItem(ctx, item.Title)
	if err != nil {
		return nil, fmt.Errorf("create item: %w", err)
	}
	return itemRowToDomain(row), nil
}

func (r *ItemRepository) Get(ctx context.Context, id int64) (*domain.Item, error) {
	row, err := r.queries.GetItem(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrItemNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}
	return itemRowToDomain(row), nil
}

func (r *ItemRepository) List(ctx context.Context) ([]*domain.Item, error) {
	rows, err := r.queries.ListItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	items := make([]*domain.Item, len(rows))
	for i, row := range rows {
		items[i] = itemRowToDomain(row)
	}
	return items, nil
}
