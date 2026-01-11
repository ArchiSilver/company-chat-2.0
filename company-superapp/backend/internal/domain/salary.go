package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Salary struct {
	UserID          uuid.UUID `db:"user_id" json:"user_id"`
	AmountEncrypted []byte    `db:"amount_encrypted" json:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type SalaryRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Salary, error)
	Upsert(ctx context.Context, salary *Salary) error
}
