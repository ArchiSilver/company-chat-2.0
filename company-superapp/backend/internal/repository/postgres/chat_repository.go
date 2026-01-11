package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/company-superapp/internal/domain"
)

type ChatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	query := `INSERT INTO messenger.chats (type) VALUES ($1) RETURNING id, created_at`
	return r.db.QueryRowxContext(ctx, query, chat.Type).Scan(&chat.ID, &chat.CreatedAt)
}

func (r *ChatRepository) AddMember(ctx context.Context, chatID, userID uuid.UUID) error {
	query := `INSERT INTO messenger.chat_members (chat_id, user_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, chatID, userID)
	return err
}

func (r *ChatRepository) GetChatsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	var chats []domain.Chat
	query := `SELECT c.id, c.type, c.created_at FROM messenger.chats c
			  JOIN messenger.chat_members cm ON c.id = cm.chat_id
			  WHERE cm.user_id = $1`
	err := r.db.SelectContext(ctx, &chats, query, userID)
	return chats, err
}
