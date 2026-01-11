package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/company-superapp/internal/domain"
)

type PushTokenRepository struct {
	db *sqlx.DB
}

func NewPushTokenRepository(db *sqlx.DB) *PushTokenRepository {
	return &PushTokenRepository{db: db}
}

func (r *PushTokenRepository) Create(ctx context.Context, token *domain.PushToken) error {
	query := `
		INSERT INTO system.push_tokens (id, user_id, token, device_info, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (token) DO UPDATE SET user_id = EXCLUDED.user_id, device_info = EXCLUDED.device_info
		RETURNING created_at
	`
	return r.db.QueryRowContext(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.DeviceInfo,
	).Scan(&token.CreatedAt)
}

func (r *PushTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.PushToken, error) {
	var tokens []domain.PushToken
	query := `SELECT id, user_id, token, device_info, created_at FROM system.push_tokens WHERE user_id = $1`
	err := r.db.SelectContext(ctx, &tokens, query, userID)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *PushTokenRepository) GetByToken(ctx context.Context, token string) (*domain.PushToken, error) {
	var pushToken domain.PushToken
	query := `SELECT id, user_id, token, device_info, created_at FROM system.push_tokens WHERE token = $1`
	err := r.db.GetContext(ctx, &pushToken, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &pushToken, nil
}

func (r *PushTokenRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM system.push_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *PushTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM system.push_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
