package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/company-superapp/internal/domain"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, msg *domain.Message) error {
	query := `INSERT INTO messenger.messages (chat_id, sender_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRowxContext(ctx, query, msg.ChatID, msg.SenderID, msg.Content).Scan(&msg.ID, &msg.CreatedAt)
}

func (r *MessageRepository) GetMessagesByChatID(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]domain.Message, error) {
	var messages []domain.Message
	query := `SELECT id, chat_id, sender_id, content, created_at FROM messenger.messages 
			  WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &messages, query, chatID, limit, offset)
	return messages, err
}
