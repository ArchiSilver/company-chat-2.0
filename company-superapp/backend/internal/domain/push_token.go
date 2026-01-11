package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PushToken struct {
	ID         uuid.UUID `db:"id" json:"id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	Token      string    `db:"token" json:"token"`
	DeviceInfo *string   `db:"device_info" json:"device_info,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type PushTokenRepository interface {
	Create(ctx context.Context, token *PushToken) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]PushToken, error)
	GetByToken(ctx context.Context, token string) (*PushToken, error)
	Delete(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
