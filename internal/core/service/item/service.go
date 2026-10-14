package item

import (
	"context"

	"github.com/1337jazz/generoni/internal/core/domain"
	"github.com/1337jazz/generoni/internal/core/port"
)

var _ port.ItemService = (*Service)(nil)

type Service struct {
	repo port.ItemRepository
}

func NewService(repo port.ItemRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, title string) (*domain.Item, error) {
	item, err := domain.NewItem(title)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, item)
}

func (s *Service) Get(ctx context.Context, id int64) (*domain.Item, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*domain.Item, error) {
	return s.repo.List(ctx)
}
