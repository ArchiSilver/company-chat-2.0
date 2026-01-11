package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/company-superapp/internal/domain"
)

type TaxiRequestRepository struct {
	db *sqlx.DB
}

func NewTaxiRequestRepository(db *sqlx.DB) *TaxiRequestRepository {
	return &TaxiRequestRepository{db: db}
}

func (r *TaxiRequestRepository) Create(ctx context.Context, request *domain.TaxiRequest) error {
	query := `
		INSERT INTO finance.taxi_requests (id, user_id, receipt_file_key, status, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING created_at
	`
	return r.db.QueryRowContext(ctx, query,
		request.ID,
		request.UserID,
		request.ReceiptFileKey,
		request.Status,
	).Scan(&request.CreatedAt)
}

func (r *TaxiRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.TaxiRequest, error) {
	var request domain.TaxiRequest
	query := `SELECT id, user_id, receipt_file_key, status, created_at FROM finance.taxi_requests WHERE id = $1`
	err := r.db.GetContext(ctx, &request, query, id)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *TaxiRequestRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.TaxiRequest, error) {
	var requests []domain.TaxiRequest
	query := `SELECT id, user_id, receipt_file_key, status, created_at FROM finance.taxi_requests WHERE user_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &requests, query, userID)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *TaxiRequestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaxiRequestStatus) error {
	query := `UPDATE finance.taxi_requests SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
