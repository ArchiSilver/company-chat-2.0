package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/company-superapp/internal/domain"
)

type SalaryRepository struct {
	db *sqlx.DB
}

func NewSalaryRepository(db *sqlx.DB) *SalaryRepository {
	return &SalaryRepository{db: db}
}

func (r *SalaryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Salary, error) {
	var salary domain.Salary
	query := `SELECT user_id, amount_encrypted, updated_at FROM finance.salaries WHERE user_id = $1`

	err := r.db.GetContext(ctx, &salary, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &salary, nil
}

func (r *SalaryRepository) Upsert(ctx context.Context, salary *domain.Salary) error {
	query := `
		INSERT INTO finance.salaries (user_id, amount_encrypted, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET amount_encrypted = EXCLUDED.amount_encrypted, updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query, salary.UserID, salary.AmountEncrypted)
	return err
}
