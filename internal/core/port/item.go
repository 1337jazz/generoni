package port

import (
	"context"

	"github.com/1337jazz/generoni/internal/core/domain"
)

type ItemService interface {
	Create(ctx context.Context, title string) (*domain.Item, error)
	Get(ctx context.Context, id int64) (*domain.Item, error)
	List(ctx context.Context) ([]*domain.Item, error)
}

type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) (*domain.Item, error)
	Get(ctx context.Context, id int64) (*domain.Item, error)
	List(ctx context.Context) ([]*domain.Item, error)
}
