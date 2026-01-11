package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        int64     `json:"id" db:"id"`
	ChatID    uuid.UUID `json:"chat_id" db:"chat_id"`
	SenderID  uuid.UUID `json:"sender_id" db:"sender_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type MessageRepository interface {
	Create(ctx context.Context, msg *Message) error
	GetMessagesByChatID(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]Message, error)
}
