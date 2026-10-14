package item

import (
	"time"

	"github.com/1337jazz/generoni/internal/core/domain"
)

type createItemRequest struct {
	Title string `json:"title"`
}

type itemResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func itemToResponse(item *domain.Item) itemResponse {
	return itemResponse{
		ID:        item.ID,
		Title:     item.Title,
		CreatedAt: item.CreatedAt,
	}
}
