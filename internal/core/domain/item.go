package domain

import (
	"strings"
	"time"
	"unicode/utf8"
)

const maxTitleLength = 200

type Item struct {
	ID        int64
	Title     string
	CreatedAt time.Time
}

func NewItem(title string) (*Item, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrTitleRequired
	}
	if utf8.RuneCountInString(title) > maxTitleLength {
		return nil, ErrTitleTooLong
	}
	return &Item{Title: title}, nil
}
