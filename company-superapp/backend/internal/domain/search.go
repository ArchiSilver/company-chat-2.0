package domain

import (
	"time"

	"github.com/google/uuid"
)

type SearchResultType string

const (
	SearchTypeUser    SearchResultType = "user"
	SearchTypeMessage SearchResultType = "message"
	SearchTypeTask    SearchResultType = "task"
)

type SearchResult struct {
	Type      SearchResultType `json:"type"`
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	Subtitle  string           `json:"subtitle,omitempty"`
	Rank      float64          `json:"rank"`
	CreatedAt time.Time        `json:"created_at"`
	// Дополнительные поля для навигации
	ChatID    *uuid.UUID `json:"chat_id,omitempty"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	TaskID    *uuid.UUID `json:"task_id,omitempty"`
	MessageID *int64     `json:"message_id,omitempty"`
}

type SearchResults struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}
