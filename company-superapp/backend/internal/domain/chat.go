package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatRepository interface {
	Create(ctx context.Context, chat *Chat) error
	AddMember(ctx context.Context, chatID, userID uuid.UUID) error
	GetChatsByUserID(ctx context.Context, userID uuid.UUID) ([]Chat, error)
}
