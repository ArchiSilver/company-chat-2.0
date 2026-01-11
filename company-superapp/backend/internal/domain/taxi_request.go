package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaxiRequestStatus string

const (
	TaxiStatusPending  TaxiRequestStatus = "pending"
	TaxiStatusApproved TaxiRequestStatus = "approved"
	TaxiStatusRejected TaxiRequestStatus = "rejected"
)

type TaxiRequest struct {
	ID             uuid.UUID         `db:"id" json:"id"`
	UserID         uuid.UUID         `db:"user_id" json:"user_id"`
	ReceiptFileKey string            `db:"receipt_file_key" json:"receipt_file_key"`
	Status         TaxiRequestStatus `db:"status" json:"status"`
	CreatedAt      time.Time         `db:"created_at" json:"created_at"`
}

type TaxiRequestRepository interface {
	Create(ctx context.Context, request *TaxiRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*TaxiRequest, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]TaxiRequest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status TaxiRequestStatus) error
}
