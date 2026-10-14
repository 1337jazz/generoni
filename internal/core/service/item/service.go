package item

import (
	"context"
	"fmt"

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

func (s *Service) List(ctx context.Context, page domain.Pagination) ([]*domain.Item, int64, error) {
	items, err := s.repo.List(ctx, page)
	if err != nil {
		return nil, 0, fmt.Errorf("list items page: %w", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count items: %w", err)
	}

	return items, total, nil
}
